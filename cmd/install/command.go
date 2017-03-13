package install

import (
	"context"
	"errors"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/self"
)

func init() {
	cmd.Register("install", Info)
}

func run(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must specify a target")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	isSelf := args[0] == "self"
	if isSelf {
		installPath := ""
		if len(args) > 1 {
			installPath = args[1]
		}

		if err := self.Install(installPath); err != nil {
			log.Printf("error installing self: %v", err)
			log.Printf("could not install locally. Depending on your system's configuration, you may need to run the install as an administrator.")
		}

		return nil
	}

	token := v1.ParseNameVersionToken(args[0])
	gamePath := manager.GetGameOrDefault(ctx.Config.Games, "")
	mod, err := manager.InstallMod(ctx, gamePath, ctx.Profile(), token)
	if err != nil {
		log.Printf("error installing mod: %v", err)
		log.Printf("could not install mod: %v", err)
		return nil
	}

	log.Printf("%s@%s installed at %s\n", mod.Name, mod.SemVersion, mod.Source)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "install",
		Short: "i",
		Long:  "install",
		Run:   cmd.ExecutorFunc(run),
	}
)
