package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielkrainas/shex/api/client"
	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/mods"
	"github.com/danielkrainas/shex/utils/sysfs"
)

const (
	DefaultGameManifestName = "shex.json"
	DefaultGameName         = "default"
	DefaultProfileName      = "default"
	HomeConfigName          = "config.json"
	HomeProfilesFolder      = "profiles"
	HomeCacheFolder         = "cache"
	HomeChannelsFolder      = "channels"
	defaultHomeFolder       = ".shex"
)

type Manager interface {
	Fs() sysfs.SysFs
	Init() error
	Config() *Config
	//Game(name string) string
	Home() string
	SaveConfig() error
	Profile() *v1.Profile
	Profiles() map[string]*v1.Profile
	Channel() *mods.Channel
	Channels() mods.ChannelMap
	AddProfile(profile *v1.Profile) error
	RemoveProfile(id string) (*v1.Profile, error)
	AddGame(alias string, game mods.GameDir) error
	RemoveGame(alias string) error
	AddChannel(ch *mods.Channel) error
	RemoveChannel(alias string) (*mods.Channel, error)
	ClearCache() error
	UninstallMod(game mods.GameDir, profile *v1.Profile, name string) error
	InstallMod(game mods.GameDir, profile *v1.Profile, token *v1.NameVersionToken) (*v1.ModInfo, error)
}

type manager struct {
	channels mods.ChannelMap
	homePath string
	profiles map[string]*v1.Profile
	config   *Config
	fs       sysfs.SysFs
}

func New(homePath string, fs sysfs.SysFs, config *Config) (Manager, error) {
	m := &manager{
		homePath: homePath,
		profiles: make(map[string]*v1.Profile),
		channels: make(mods.ChannelMap),
		fs:       fs,
		config:   config,
	}

	if err := m.Init(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *manager) Init() error {
	if err := m.loadProfiles(); err != nil {
		return err
	}

	if err := m.loadChannels(); err != nil {
		return err
	}

	if m.config.IncludeDefaultChannel {
		m.channels[DefaultChannel.Alias] = DefaultChannel
	}

	return nil
}

func (m *manager) Config() *Config {
	return m.config
}

func (m *manager) Home() string {
	return m.homePath
}

func (m *manager) Fs() sysfs.SysFs {
	return m.fs
}

func (m *manager) Profile() *v1.Profile {
	return m.profiles[m.config.ActiveProfile]
}

func (m *manager) Profiles() map[string]*v1.Profile {
	return m.profiles
}

func (m *manager) AddProfile(profile *v1.Profile) error {
	if _, ok := m.profiles[profile.Id]; ok {
		return fmt.Errorf("[%s] already exists", profile.Id)
	}

	m.profiles[profile.Id] = profile
	return nil
}

func (m *manager) Channel() *mods.Channel {
	return m.channels[m.config.ActiveRemote]
}

func (m *manager) Channels() mods.ChannelMap {
	return m.channels
}

func (m *manager) UninstallMod(game mods.GameDir, profile *v1.Profile, name string) error {
	gameManifest, err := mods.LoadGameManifest(m.fs, game.String())
	if err != nil {
		return err
	}

	profile.Mods.Drop(name)
	if err := m.saveProfile(profile); err != nil {
		return err
	}

	if err := mods.Uninstall(m.fs, game, gameManifest, name); err != nil {
		return err
	}

	gameManifest.Mods.Drop(name)
	if err := mods.SaveGameManifest(game.String(), gameManifest); err != nil {
		return err
	}

	return nil
}

func (m *manager) InstallMod(game mods.GameDir, profile *v1.Profile, token *v1.NameVersionToken) (*v1.ModInfo, error) {
	// keep moving to sysFs and mods provider
	ch := m.Channel()
	source := ch.Protocol + "://" + ch.Endpoint
	remoteInfo, err := client.DownloadModInfo(source, token)
	if err != nil {
		return nil, err
	}

	localName := mods.GetLocalModPathName(remoteInfo.Name, remoteInfo.Version)
	localPath := filepath.Join(game.String(), mods.ModsFolder, localName)
	err = client.DownloadMod(source, localPath, remoteInfo)
	if err != nil {
		return nil, err
	}

	gameManifest, err := mods.LoadGameManifest(m.fs, game.String())
	if err != nil {
		return nil, err
	}

	profile.Mods.Set(remoteInfo.Name, token.Constraint())
	if err := m.saveProfile(profile); err != nil {
		return nil, err
	}

	gameManifest.Mods.Set(remoteInfo.Name, remoteInfo.Version)
	err = mods.SaveGameManifest(game.String(), gameManifest)
	if err != nil {
		return nil, err
	}

	return mods.GetModInfo(m.fs, localPath)
}

func (m *manager) RemoveProfile(id string) (*v1.Profile, error) {
	p, ok := m.profiles[id]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", id)
	}

	if err := m.dropProfile(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (m *manager) dropProfile(profile *v1.Profile) error {
	profilePath := m.pathFor(profile)
	if !m.fs.FileExists(profilePath) {
		return nil // TODO return error here with message?
	}

	if err := m.fs.DeleteFile(profilePath); err != nil {
		return err
	}

	return nil
}

func (m *manager) saveProfile(profile *v1.Profile) error {
	return SaveProfile(m.Fs(), m.pathFor(profile), profile)
}

func (m *manager) loadProfiles() error {
	files, err := m.fs.ReadDir(m.config.ProfilesPath)
	if err != nil {
		return err
	}

	result := make(map[string]*v1.Profile)
	for _, f := range files {
		if isJson, err := filepath.Match("*.json", f.Name()); err != nil {
			return err
		} else if !isJson {
			continue
		}

		if profile, err := LoadProfile(m.fs, filepath.Join(m.config.ProfilesPath, f.Name())); err != nil {
			return err
		} else {
			result[profile.Id] = profile
		}
	}

	m.profiles = result
	return nil
}

func (m *manager) dropChannel(ch *mods.Channel) error {
	channelPath := m.pathFor(ch)
	if !m.fs.FileExists(channelPath) {
		return nil // TODO return error here with message?
	}

	if err := m.fs.DeleteFile(channelPath); err != nil {
		return err
	}

	return nil
}

func (m *manager) saveChannel(ch *mods.Channel) error {
	return sysfs.WriteJson(m.fs, m.pathFor(ch), ch)
}

func (m *manager) loadChannel(channelPath string) (*mods.Channel, error) {
	var ch *mods.Channel
	if err := sysfs.ReadJson(m.fs, channelPath, ch); err != nil {
		return nil, err
	}

	return ch, nil
}

func (m *manager) loadChannels() error {
	files, err := m.fs.ReadDir(m.config.ChannelsPath)
	if err != nil {
		return err
	}

	result := make(mods.ChannelMap)
	for _, f := range files {
		if isJson, err := filepath.Match("*.json", f.Name()); err != nil {
			return err
		} else if !isJson {
			continue
		}

		if channel, err := m.loadChannel(filepath.Join(m.config.ChannelsPath, f.Name())); err != nil {
			return err
		} else {
			result[channel.Alias] = channel
		}
	}

	m.channels = result
	return nil
}

func (m *manager) SaveConfig() error {
	return SaveConfig(m.fs, m.homePath, m.config)

}

func (m *manager) AddGame(alias string, dir mods.GameDir) error {
	m.config.Games.Attach(alias, dir)
	if err := SaveConfig(m.fs, m.homePath, m.config); err != nil {
		return err
	}

	return nil
}

func (m *manager) RemoveGame(alias string) error {
	_, ok := m.config.Games[alias]
	if !ok {
		fmt.Printf("game %q does not exist.", alias)
		return nil
	}

	m.config.Games.Detach(alias)
	return nil
}

func (m *manager) AddChannel(ch *mods.Channel) error {
	m.channels[ch.Alias] = ch
	return m.saveChannel(ch)
}

func (m *manager) RemoveChannel(alias string) (*mods.Channel, error) {
	var ch *mods.Channel
	ok := false
	if alias == "default" && m.config.IncludeDefaultChannel {
		ch = DefaultChannel
		ok = true
	} else {
		ch, ok = m.channels[alias]
	}

	if !ok {
		return nil, fmt.Errorf("channel %q not found", alias)
	}

	if ch == DefaultChannel {
		m.config.IncludeDefaultChannel = false
	} else if err := m.dropChannel(ch); err != nil {
		return nil, err
	}

	return ch, nil
}

func (m *manager) pathFor(v interface{}) string {
	switch t := v.(type) {
	case *v1.Profile:
		return filepath.Join(m.config.ProfilesPath, t.Id+".json")
	case *mods.Channel:
		return filepath.Join(m.config.ChannelsPath, t.Alias+".json")
	case string:
		if t == "cache" {
			return filepath.Join(m.homePath, m.config.CachePath)
		}
	}

	return ""
}

func (m *manager) ClearCache() error {
	return sysfs.ClearDir(m.fs, m.pathFor("cache"))
}

func GetGameOrDefault(games mods.GameMap, name string) mods.GameDir {
	if name == "" {
		name = DefaultGameName
	}

	game, ok := games[name]
	if !ok {
		return mods.EmptyGame
	}

	return game
}

func getDefaultHomePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(u.HomeDir, string(filepath.Separator)+defaultHomeFolder), nil
}

func ensureHomeDirectoryExists(fs sysfs.SysFs, homePath string) error {
	if !fs.DirExists(homePath) {
		err := os.Mkdir(homePath, 0777)
		if err != nil {
			return err
		}
	}

	configPath := filepath.Join(homePath, HomeConfigName)
	profilesPath := filepath.Join(homePath, HomeProfilesFolder)
	cachePath := filepath.Join(homePath, HomeCacheFolder)
	channelsPath := filepath.Join(homePath, HomeChannelsFolder)
	if !fs.FileExists(configPath) {
		defaultConfig := NewConfig()
		defaultConfig.ProfilesPath = profilesPath

		jsonContent, err := json.Marshal(&defaultConfig)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(configPath, jsonContent, 0777); err != nil {
			return err
		}
	}

	if !fs.DirExists(cachePath) {
		if err := fs.CreateDir(cachePath); err != nil {
			return err
		}
	}

	if !fs.DirExists(channelsPath) {
		if err := fs.CreateDir(channelsPath); err != nil {
			return err
		}
	}

	if !fs.DirExists(profilesPath) {
		if err := fs.CreateDir(profilesPath); err != nil {
			return err
		}
	}

	defaultProfilePath := path.Join(profilesPath, DefaultProfileName+".json")
	if !fs.FileExists(defaultProfilePath) {
		defaultProfile := v1.Profile{}
		defaultProfile.Id = DefaultProfileName
		defaultProfile.Mods = make(map[string]string)
		defaultProfile.Name = strings.ToTitle(DefaultProfileName)

		jsonContent, err := json.Marshal(&defaultProfile)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(defaultProfilePath, jsonContent, 0777); err != nil {
			return err
		}
	}

	return nil
}
