package mods

import (
	"errors"
	"path"
	"regexp"
	"strings"

	"github.com/danielkrainas/shex/utils/sysfs"
)

const DefaultGameManifestName = "shex.json"

const (
	ModsFolder         = "mods"
	EmptyGame  GameDir = ""
)

type GameDir string

func (dir GameDir) String() string {
	return string(dir)
}

type GameMap map[string]GameDir

func (games GameMap) Len() int {
	return len(games)
}

func (games GameMap) Contains(name string) bool {
	_, ok := games[name]
	return ok
}

func (games GameMap) Attach(name string, dir GameDir) {
	games[name] = dir
}

func (games GameMap) Detach(name string) {
	delete(games, name)
}

func getModBasePath(modPath string) string {
	base := path.Base(modPath)
	return strings.TrimRight(base, ".smod")
}

func FindGameVersion(fs sysfs.SysFs) (string, error) {
	versionRegEx := regexp.MustCompile("((?:Develop|(?:[a-zA-Z]+))\\-[0-9]+)")
	notesPath := "stonehearth/release_notes/release_notes.html"
	rawContent, err := getZipResourceContent(fs, "/home/daniel/Documents/stonehearth.smod", notesPath, false)
	result := ""
	if err == nil {
		if len(rawContent) <= 0 {
			return "", errors.New("could not find release notes")
		}

		notesContent := string(rawContent[:])
		for _, line := range strings.Split(notesContent, "\n") {
			if strings.Contains(line, "<h2>") && versionRegEx.MatchString(line) {
				result = versionRegEx.FindString(line)
				break
			}
		}
	}

	return result, err
}
