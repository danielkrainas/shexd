package main

import (
	"context"
	"math/rand"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	_ "github.com/danielkrainas/shexd/blobs/driver/inmemory"
	_ "github.com/danielkrainas/shexd/cmd/api"
	"github.com/danielkrainas/shexd/cmd/root"
	_ "github.com/danielkrainas/shexd/cmd/version"
	_ "github.com/danielkrainas/shexd/storage/driver/inmemory"
)

var appVersion string

const defaultVersion = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = defaultVersion
	}

	rand.Seed(time.Now().Unix())
	ctx := acontext.WithVersion(acontext.Background(), appVersion)
	ctx = context.WithValue(ctx, "app.name", strings.Title(root.Info.Use))

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
