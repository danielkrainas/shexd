package use

import (
	"context"
	"errors"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("use", Info)
}

func run(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("profile name not specified")
	}

	mctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	newProfileName := args[0]
	if newProfileName != mctx.Config.ActiveProfile {
		newProfile := mctx.Profiles[newProfileName]
		mctx.Config.ActiveProfile = newProfile.Id
		if err := manager.SaveConfig(mctx.Config, mctx.HomePath); err != nil {
			return err
		}

		log.Printf("active profile set to: %s\n", newProfile.Name)
	} else {
		log.Printf("profile already active")
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "use",
		Short: "use",
		Long:  "use",
		Run:   cmd.ExecutorFunc(run),
	}
)
