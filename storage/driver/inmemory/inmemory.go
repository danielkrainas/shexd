package inmemory

import (
	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/shexd/api/v1"
	"github.com/danielkrainas/shexd/storage"
	"github.com/danielkrainas/shexd/storage/driver/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (drivers.DriverBase, error) {
	return &driver{
		stores: make(map[string]interface{}, 0),
	}, nil
}

func init() {
	factory.Register("inmemory", &driverFactory{})
}

const (
	modStoreKey     = "mods"
	profileStoreKey = "profiles"
)

type driver struct {
	stores map[string]interface{}
}

func (d *driver) Mods() storage.ModStore {
	store, ok := d.stores[modStoreKey].(storage.ModStore)
	if !ok {
		store = &modStore{
			mods: make([]*v1.ModInfo, 0),
		}

		d.stores[modStoreKey] = store
	}

	return store
}

func (d *driver) Profiles() storage.ProfileStore {
	store, ok := d.stores[profileStoreKey].(storage.ProfileStore)
	if !ok {
		store = &profileStore{
			profiles: make([]*v1.RemoteProfile, 0),
		}

		d.stores[profileStoreKey] = store
	}

	return store
}
