package version

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"
)

func init() {
	cmd.Register("version", Info)
}

func run(ctx context.Context, args []string) error {
	fmt.Printf("%s v%s\n", acontext.GetStringValue(ctx, "app.name"), acontext.GetVersion(ctx))
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "version",
		Short: "show version information",
		Long:  "show version information",
		Run:   cmd.ExecutorFunc(run),
	}
)
