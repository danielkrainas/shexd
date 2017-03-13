package main

import (
	"context"
	"math/rand"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	_ "github.com/danielkrainas/shex/cmd/api"
	"github.com/danielkrainas/shex/cmd/registry"
	_ "github.com/danielkrainas/shex/cmd/version"
	_ "github.com/danielkrainas/shex/registry/storage/driver/inmemory"
)

var appVersion string

const defaultVersion = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = defaultVersion
	}

	rand.Seed(time.Now().Unix())
	ctx := acontext.WithVersion(acontext.Background(), appVersion)
	ctx = context.WithValue(ctx, "app.name", strings.Title(registry.Info.Use))

	dispatch := cmd.CreateDispatcher(ctx, registry.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
