package storage

import (
	"errors"

	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/shex/api/v1"
)

var ErrNotFound = errors.New("not found")

type Driver interface {
	drivers.DriverBase

	Mods() ModStore
}

type ModStore interface {
	Delete(token *v1.NameVersionToken) error
	Store(m *v1.ModInfo, isNew bool) error
	Find(token *v1.NameVersionToken) (*v1.ModInfo, error)
	FindMany(f *ModFilters) ([]*v1.ModInfo, error)
	Count(f *ModFilters) (int, error)
}

type ModFilters struct{}
