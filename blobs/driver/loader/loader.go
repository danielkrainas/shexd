package loader

import (
	"errors"

	"github.com/danielkrainas/shexd/configuration"
	"github.com/danielkrainas/shexd/driver"
)

func FromConfig(config *configuration.Config) (driver.Driver, error) {
	params := config.Blobs.Parameters()
	if params == nil {
		params = make(configuration.Parameters)
	}

	d, err := factory.Create(config.Blobs.Type(), params)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func LogSummary(ctx context.Context, config *configuration.Config) {
	acontext.GetLogger(ctx).Infof("using %q blobs driver", config.Blobs.Type())
}
