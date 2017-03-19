package factory

import (
	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/shexd/blobs/driver"
)

var registry = &drivers.Registry{
	AssetType: "Blob Storage",
}

func Register(name string, factory drivers.Factory) {
	registry.Register(name, factory)
}

func Create(name string, parameters map[string]interface{}) (driver.Driver, error) {
	d, err := registry.Create(name, parameters)
	return d.(driver.Driver), err
}
