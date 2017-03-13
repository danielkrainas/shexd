package manager

import (
	"path/filepath"

	"github.com/danielkrainas/shex/mods"
	"github.com/danielkrainas/shex/utils/sysfs"
)

type Config struct {
	filePath              string
	ActiveProfile         string       `json:"active"`
	ActiveRemote          string       `json:"remote"`
	ProfilesPath          string       `json:"profiles"`
	ChannelsPath          string       `json:"channels"`
	IncludeDefaultChannel bool         `json:"includeDefaultChannel"`
	CachePath             string       `json:"cache"`
	Games                 mods.GameMap `json:"games"`
}

func SaveConfig(fs sysfs.SysFs, homePath string, config *Config) error {
	var err error
	if len(homePath) <= 0 {
		homePath, err = getDefaultHomePath()
		if err != nil {
			return err
		}
	}

	configPath := filepath.Join(homePath, HomeConfigName)
	if err := sysfs.WriteJson(fs, configPath, config); err != nil {
		return err
	}

	return nil
}

func NewConfig() *Config {
	return &Config{
		ActiveProfile: DefaultProfileName,
		ActiveRemote:  "default",
		Games:         make(mods.GameMap),
		IncludeDefaultChannel: true,
	}
}

func LoadConfig(fs sysfs.SysFs, homePath string) (*Config, error) {
	config := NewConfig()
	var err error
	if homePath == "" {
		homePath, err = getDefaultHomePath()
		if err != nil {
			return nil, err
		}
	}

	if err = ensureHomeDirectoryExists(fs, homePath); err != nil {
		return nil, err
	}

	configPath := filepath.Join(homePath, HomeConfigName)
	if err := sysfs.ReadJson(fs, configPath, config); err != nil {
		return nil, err
	}

	if len(config.ProfilesPath) < 1 {
		config.ProfilesPath = filepath.Join(homePath, HomeProfilesFolder)
	}

	if len(config.ChannelsPath) < 1 {
		config.ChannelsPath = filepath.Join(homePath, HomeChannelsFolder)
	}

	return config, nil
}
