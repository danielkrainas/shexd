package queries

import (
	"github.com/danielkrainas/shex/api/v1"
)

type CountMods struct{}

type FindMod struct {
	Token *v1.NameVersionToken
}

type SearchMods struct{}
