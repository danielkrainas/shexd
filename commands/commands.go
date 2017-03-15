package commands

import (
	"github.com/danielkrainas/shexd/api/v1"
)

type DeleteMod struct {
	Token *v1.NameVersionToken
}

type StoreMod struct {
	New bool
	Mod *v1.ModInfo
}

type StoreProfile struct {
	New     bool
	Profile *v1.RemoteProfile
}
