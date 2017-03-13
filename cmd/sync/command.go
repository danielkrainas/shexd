package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("sync", Info)
}

var (
	Info = &cmd.Info{
		Use:   "sync",
		Short: "sync",
		Long:  "sync",
		SubCommands: []*cmd.Info{
			{
				Use:   "profile",
				Short: "profile",
				Long:  "profile",
				Run:   cmd.ExecutorFunc(syncProfile),
			},
			{
				Use:   "profiles",
				Short: "profiles",
				Long:  "profiles",
				Run:   cmd.ExecutorFunc(syncAllProfiles),
			},
		},
	}
)

func reportSyncResult(artifactName string, fromVersion string, toVersion string) {
	if fromVersion == toVersion {
		fmt.Printf("%s %-20s\n", artifactName, "OK")
	} else {
		fmt.Printf("%s %-20s->%s\n", artifactName, fromVersion, toVersion)
	}
}

func reportProfileSyncResult(p *v1.Profile, from int32, to int32) {
	if from == to {
		fmt.Printf("%s @%d => no updates available\n", p.Name, from)
	} else {
		fmt.Printf("%s @%d => @%d\n", p.Name, from, to)
	}
}

/* Sync Profiles Command */
func syncAllProfiles(parent context.Context, args []string) error {
	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	for _, p := range ctx.Profiles {
		if p.Source == nil {
			continue
		}

		from, to, err := manager.SyncProfile(p)
		if err != nil {
			log.Println("ERR: error syncing profile: %v", err)
			log.Println("couldn't sync with remote server.")
			return nil
		}

		err = manager.SaveProfile(p)
		if err != nil {
			log.Printf("error saving profile: %v", err)
			log.Println("couldn't save profile")
			return nil
		}

		reportProfileSyncResult(p, from, to)
	}

	return nil
}

/* Sync Profile Command */
func syncProfile(parent context.Context, args []string) error {
	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	profile := ctx.Profile()
	if len(args) > 0 {
		var ok bool
		profile, ok = ctx.Profiles[args[0]]
		if !ok {
			log.Println("profile not found.")
			return nil
		}
	}

	if profile.Source == nil {
		log.Println("not a remote profile")
		return nil
	}

	from, to, err := manager.SyncProfile(profile)
	if err != nil {
		log.Printf("error syncing profile: %v", err)
		log.Println("couldn't sync with remote server.")
		return nil
	}

	if err = manager.SaveProfile(profile); err != nil {
		log.Printf("error saving profile: %v", err)
		log.Println("couldn't save profile")
		return nil
	}

	reportProfileSyncResult(profile, from, to)
	return nil
}
