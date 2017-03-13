package cmdutils

import (
	"context"

	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shexd/manager"
	"github.com/danielkrainas/shexd/utils/sysfs"
)

func LoadManager(ctx context.Context) (manager.Manager, error) {
	homePath := acontext.GetStringValue(ctx, "flags.home")
	//configPath := acontext.GetStringValue(ctx, "flags.config")

	fs := sysfs.New()
	config, err := manager.LoadConfig(fs, homePath)
	if err != nil {
		return nil, err
	}

	return manager.New(homePath, fs, config)
}
