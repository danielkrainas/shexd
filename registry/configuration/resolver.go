package configuration

import (
	"fmt"
	"os"
)

func Resolve(args []string) (*Config, error) {
	var configPath string

	if len(args) > 0 {
		configPath = args[0]
	} else {
		configPath = os.Getenv("SHEXR_CONFIG_PATH")
	}

	if configPath == "" {
		return nil, fmt.Errorf("configuration path not specified")
	}

	fp, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	config, err := Parse(fp)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", configPath, err)
	}

	return config, nil
}
