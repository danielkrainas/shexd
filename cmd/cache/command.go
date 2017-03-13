package cache

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/cmd/cmdutils"
)

func init() {
	cmd.Register("cache", Info)
}

var (
	Info = &cmd.Info{
		Use:   "cache",
		Short: "perform operations on the local cache",
		Long:  "Perform operations on the local cache.",
		SubCommands: []*cmd.Info{
			{
				Use:   "clean",
				Short: "clears the local cache",
				Long:  "Clears the local cache of all contents.",
				Run:   cmd.ExecutorFunc(cleanCache),
			},
		},
	}
)

func cleanCache(ctx context.Context, _ []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if err := m.ClearCache(); err != nil {
		fmt.Printf("error clearing cache: %v\n", err)
		return nil
	}

	fmt.Println("cache cleared")
	return nil
}
