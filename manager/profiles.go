package manager

import (
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/danielkrainas/shex/api/client"
	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/utils/sysfs"
)

func SyncProfile(p *v1.Profile) (int32, int32, error) {
	if p.Source == nil {
		return 0, 0, nil
	}

	rp, err := client.DownloadProfile(p.Source)
	if err != nil {
		return 0, 0, err
	}

	old := p.Revision
	p.Mods = rp.Mods
	p.Revision = rp.Revision
	return old, p.Revision, nil
}

func pullProfile(source *v1.ProfileSource, localName string, profilesPath string) (*v1.Profile, error) {
	if source.Type != "remote" {
		return nil, errors.New("source type not supported")
	}

	profile, err := client.DownloadProfileAsLocal(source, localName)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func pushProfile(profile *v1.Profile, remoteName string, endpoint string) (string, error) {
	url := endpoint + "profiles/" + remoteName
	remoteProfile := *profile
	remoteProfile.Id = remoteName
	remoteProfile.Name = strings.Title(path.Base(remoteName))
	remoteProfile.Source = nil
	remoteProfile.Revision = 0
	jsonContent, err := json.Marshal(remoteProfile)
	if err != nil {
		return "", err
	}

	res, err := client.PostContent(url, jsonContent)
	if err != nil {
		return "", err
	}

	return string(res[:]), nil
}

func createProfileSource(name string, location string) v1.ProfileSource {
	source := v1.ProfileSource{}
	source.Location = location
	source.Uid = name
	source.Type = "remote"
	return source
}

func LoadProfile(fs sysfs.SysFs, profilePath string) (*v1.Profile, error) {
	profile := &v1.Profile{}
	if err := sysfs.ReadJson(fs, profilePath, profile); err != nil {
		return nil, err
	}

	if profile.Source != nil && profile.Source.Type == v1.SOURCE_NONE {
		profile.Source = nil
	}

	return profile, nil
}

func SaveProfile(fs sysfs.SysFs, filePath string, profile *v1.Profile) error {
	return sysfs.WriteJson(fs, filePath, profile)
}
