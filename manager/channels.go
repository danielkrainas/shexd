package manager

import (
	"github.com/danielkrainas/shex/mods"
)

var (
	DefaultChannel = &mods.Channel{
		Alias:    "default",
		Endpoint: "127.0.0.1:6231/",
		Protocol: "http",
	}
)
