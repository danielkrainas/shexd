package queries

import (
	"github.com/danielkrainas/shexd/api/v1"
)

type CountMods struct{}

type FindMod struct {
	Token *v1.NameVersionToken
}

type GetModVersionList struct {
	Token *v1.NameVersionToken
}

type SearchMods struct{}

type SearchProfiles struct{}
