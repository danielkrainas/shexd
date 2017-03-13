package uninstall

import (
	"context"
	"errors"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/self"
)

func init() {
	cmd.Register("uninstall", Info)
}

func run(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("profile name not specified")
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

		if err := self.Uninstall(installPath); err != nil {
			log.Printf("error uninstalling self: %v", err)
			log.Println("Could not uninstall. Depending on your system's configuration, you may need to run the uninstall again as an administrator.")
		}

		return nil
	}

	name := args[0]
	gamePath := manager.GetGameOrDefault(ctx.Config.Games, name)
	mod, err := manager.UninstallMod(ctx.Config, gamePath, ctx.Profile(), name)
	if err != nil {
		log.Printf("error uninstalling mod: %v", err)
		log.Println("Could not uninstall mod")
		return nil
	}

	log.Printf("%s@%s uninstalled", mod.Name, mod.SemVersion)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "uninstall",
		Short: "u",
		Long:  "uninstall",
		Run:   cmd.ExecutorFunc(run),
	}
)
