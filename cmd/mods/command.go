package mods

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/cmd/cmdutils"
	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/mods"
)

func init() {
	cmd.Register("mods", Info)
}

var (
	Info = &cmd.Info{
		Use:   "mods",
		Short: "mods",
		Long:  "mods",
		SubCommands: []*cmd.Info{
			{
				Use:   "list",
				Short: "list",
				Long:  "list",
				Run:   cmd.ExecutorFunc(listMods),
				Flags: []*cmd.Flag{
					{
						Long:        "profile",
						Short:       "p",
						Description: "display mods installed in a profile",
					},
				},
			},
		},
	}
)

/* List Mods Command */
func listMods(ctx context.Context, args []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	profileName := acontext.GetStringValue(ctx, "flags.profile")
	useProfile := profileName != ""
	var list v1.ModList
	if useProfile {
		if len(profileName) > 0 {
			selectedProfile, ok := m.Profiles()[profileName]
			if !ok {
				return fmt.Errorf("profile not found: %q", profileName)
			}

			list = selectedProfile.Mods
		} else {
			profileName = m.Profile().Name
			list = m.Profile().Mods
		}
	} else if len(m.Config().Games) <= 0 {
		fmt.Println("no games attached")
		return nil
	} else {
		gameName := ""
		if len(args) > 0 {
			gameName = args[0]
		}

		game := manager.GetGameOrDefault(m.Config().Games, gameName)
		manifest, err := mods.LoadGameManifest(m.Fs(), game.String())
		if err != nil {
			fmt.Printf("error loading game manifest: %v", err)
			fmt.Println("game manifest not found or invalid")
			return nil
		}

		list = manifest.Mods
	}

	//fmt.Printf("%-30s   %s\n", "NAME", "VERSION")
	if len(list) > 0 {
		if useProfile {
			fmt.Printf("Mods installed in profile %s\n", profileName)
		}

		for name, version := range list {
			fmt.Printf("%15s@%s\n", name, version)
		}
	} else {
		fmt.Printf("no mods installed\n")
	}

	return nil
}
