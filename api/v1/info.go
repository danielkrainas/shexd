package v1

import (
	"strings"
)

const (
	SOURCE_REMOTE = "remote"
	SOURCE_NONE   = ""
)

type ProfileSource struct {
	Type     string `json:"type"`
	Uid      string `json:"uid"`
	Location string `json:"url"`
}

type RemoteProfile struct {
	Source   *ProfileSource `json:"-"`
	Name     string         `json:"name"`
	Mods     ModList        `json:"mods"`
	Revision int32          `json:"rev"`
}

type Profile struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Mods     ModList        `json:"mods"`
	Source   *ProfileSource `json:"source"`
	Revision int32          `json:"rev"`
}

func NewProfile(id string) *Profile {
	profile := &Profile{
		Id:       id,
		Name:     strings.Title(id),
		Mods:     make(ModList),
		Revision: 1,
		Source:   nil,
	}

	return profile
}

func NewRemoteProfile(source *ProfileSource) *RemoteProfile {
	profile := &RemoteProfile{
		Mods:   make(ModList),
		Source: source,
	}

	return profile
}

func MakeLocalProfile(localName string, remote *RemoteProfile) *Profile {
	profile := NewProfile(localName)
	profile.Source = remote.Source
	profile.Revision = remote.Revision
	profile.Mods = remote.Mods
	return profile
}

type RemoteModInfo struct {
	Source  string `json:"-"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ModInfo struct {
	Source     string `json:"-"`
	Name       string `json:"name"`
	Version    int32  `json:"version"`
	SemVersion string `json:"semversion"`
}

type ModList map[string]string

func (l ModList) Contains(modName string) bool {
	_, ok := l[modName]
	return ok
}

func (l ModList) Set(name string, version string) {
	l[name] = version
}

func (l ModList) Drop(name string) {
	delete(l, name)
}

type NameVersionToken struct {
	Name    string
	Version string
}

func (t *NameVersionToken) Constraint() string {
	if t.Version != "latest" {
		return "^" + t.Version
	}

	return t.Version
}

func ParseNameVersionToken(pair string) *NameVersionToken {
	token := &NameVersionToken{}
	parts := strings.Split(pair, "@")
	token.Name = parts[0]
	if len(parts) > 1 {
		token.Version = parts[1]
	} else {
		token.Version = "latest"
	}

	return token
}
