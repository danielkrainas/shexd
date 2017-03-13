package mods

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielkrainas/shex/api/client"
	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/utils/sysfs"
)

type ModManifest struct {
	Info v1.ModInfo `json:"info"`
}

type GameManifest struct {
	Mods v1.ModList `json:"mods"`
}

func GetLocalModPathName(remoteName string, version string) string {
	parts := strings.Split(remoteName, "/")
	return parts[0] + "_" + parts[1] + "-" + version + ".smod"
}

/*func isVersionOkay(version string, versionConstraint string) bool {

}
*/

func CreateGameManifest() *GameManifest {
	manifest := &GameManifest{
		Mods: make(v1.ModList),
	}

	return manifest
}

func LoadGameManifest(fs sysfs.SysFs, gamePath string) (*GameManifest, error) {
	manifest := CreateGameManifest()
	manifestPath := filepath.Join(gamePath, DefaultGameManifestName)
	if !fs.FileExists(manifestPath) {
		return manifest, nil
	}

	if err := sysfs.ReadJson(fs, manifestPath, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func SaveGameManifest(gamePath string, manifest *GameManifest) error {
	jsonContent, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	manifestPath := path.Join(gamePath, DefaultGameManifestName)
	return ioutil.WriteFile(manifestPath, jsonContent, 0777)
}

func getZipResourceContent(fs sysfs.SysFs, zipPath string, resourcePath string, isRelative bool) ([]byte, error) {
	fp, err := fs.Read(zipPath)
	if err != nil {
		return nil, err
	}

	size, err := sysfs.SizeOf(fp)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	zr, err := zip.NewReader(fp.(io.ReaderAt), size)
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		if isRelative {
			resourcePath = path.Join(strings.Split(f.Name, "/")[0], resourcePath)
			isRelative = false
		}
		//keep switching to use SysFs
		if f.Name == resourcePath {
			r, err := f.Open()
			contents, err := ioutil.ReadAll(r)
			r.Close()
			return contents, err
		}
	}

	return nil, err
}

func Uninstall(fs sysfs.SysFs, game GameDir, manifest *GameManifest, name string) error {
	modsPath := filepath.Join(game.String(), ModsFolder)
	version, ok := manifest.Mods[name]
	if ok {
		modPath := filepath.Join(modsPath, getLocalModPathName(name, version))
		if fs.FileExists(modPath) {
			if err := fs.DeleteFile(modPath); err != nil {
				return err
			}
		}
	} else {
		//fmt.Printf("not installed in %s: \"%s\"\n", modsPath, name)
	}

	return nil
}

func Install(fs sysfs.SysFs, game GameDir, ch *Channel, profile *v1.Profile, token *v1.NameVersionToken) (*v1.ModInfo, error) {
	source := ch.Protocol + "://" + ch.Endpoint
	remoteInfo, err := client.DownloadModInfo(source, token)
	if err != nil {
		return nil, err
	}

	localName := getLocalModPathName(remoteInfo.Name, remoteInfo.Version)
	localPath := filepath.Join(game.String(), ModsFolder, localName)
	if err := client.DownloadMod(source, localPath, remoteInfo); err != nil {
		return nil, err
	}

	return GetModInfo(fs, localPath)
}

func GetModInfo(fs sysfs.SysFs, modPath string) (*v1.ModInfo, error) {
	manifestPath := "/manifest.json"
	jsonContent, err := getZipResourceContent(fs, modPath, manifestPath, true)
	if err != nil {
		return nil, err
	}

	if len(jsonContent) <= 0 {
		return nil, errors.New("could not find manifest file")
	}

	manifest := ModManifest{}
	err = json.Unmarshal(jsonContent, &manifest)
	if err == nil {
		if len(manifest.Info.SemVersion) <= 0 {
			manifest.Info.SemVersion = fmt.Sprintf("%d.0.0", manifest.Info.Version)
		}

		manifest.Info.Source = modPath
	}

	return &manifest.Info, err
}

func getLocalModPathName(remoteName string, version string) string {
	parts := strings.Split(remoteName, "/")
	return parts[0] + "_" + parts[1] + "-" + version + ".smod"
}

/*func isVersionOkay(version string, versionConstraint string) bool {

}

func isModCached(config *ManagerConfig) bool {

}*/

/*func execStat(current *executionContext) error {
	modPath := args[0]
	info, err := getModInfo(modPath)
	if err != nil {
		return appError{err, "Could not find mod information"}
	}

	log.Printf("[%s]\n", modPath)
	log.Printf("name: %s\nversion: %d\nsem version: %s\n", info.Name, info.Version, info.SemVersion)
	return nil
}

func execPull(current *executionContext) error {
	remoteName := args[0]
	localName := path.Base(remoteName)
	if len(current.args) > 1 {
		localName = args[1]
	}

	if _, ok := current.profiles[localName]; ok {
		return appError{nil, fmt.Sprintf("[%s] already exists", localName)}
	}

	var ok bool
	// TODO: put this together as its own part later
	current.remote = getDefaultRemote()
	if current.config.ActiveRemote != DefaultRemoteName {
		current.remote, ok = current.config.Remotes[current.config.ActiveRemote]
		if !ok {
			return appError{nil, fmt.Sprintf("remote \"%s\" not found\n", current.config.ActiveRemote)}
		}
	}

	source := createProfileSource(remoteName, current.remote)
	profile, err := pullProfile(&source, localName, current.config.ProfilesPath)
	if err != nil {
		return appError{err, "Could not pull profile from the server"}
	}

	log.Printf("pulled [%s] to: %s\n", profile.Source.Uid, profile.filePath)
	return nil
}

func execPush(current *executionContext) error {
	profileId := args[0]
	profile, ok := current.profiles[profileId]
	if !ok {
		return appError{nil, fmt.Sprintf("[%s] not found\n", profileId)}
	}

	remoteName := args[1]
	if remoteName != current.config.ActiveRemote {
		current.remote = getRemoteOrDefault(current.config.Remotes, remoteName)
	}

	// TODO: pull out core logic into own func or something
	remote := getDefaultRemote()
	if current.config.ActiveRemote != "default" {
		remote, ok = current.config.Remotes[current.config.ActiveRemote]
		if !ok {
			return appError{nil, fmt.Sprintf("remote \"%s\" not found\n", current.config.ActiveRemote)}
		}
	}

	version, err := pushProfile(profile, remoteName, remote)
	if err != nil {
		return appError{err, "Could not push profile to server"}
	}

	log.Printf("[%s] pushed to %s as %s@%s\n", profileId, current.config.ActiveRemote, remoteName, version)
	log.Printf("import with: `shex pull %s@%s`\n", remoteName, version)
	return nil
}*/
