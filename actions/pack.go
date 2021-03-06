package actions

import (
	"context"

	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/shexd/commands"
	"github.com/danielkrainas/shexd/configuration"
	"github.com/danielkrainas/shexd/queries"
	"github.com/danielkrainas/shexd/storage"
	"github.com/danielkrainas/shexd/storage/loader"
)

type Pack interface {
	cqrs.QueryExecutor
	cqrs.CommandHandler
}

type pack struct {
	store storage.Driver
}

func (p *pack) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.FindMod:
		return FindMod(ctx, q, p.store.Mods())
	case *queries.GetModVersionList:
		return GetModVersionList(ctx, q, p.store.Mods())
	case *queries.CountMods:
		return CountMods(ctx, q, p.store.Mods())
	case *queries.SearchMods:
		return SearchMods(ctx, q, p.store.Mods())
	case *queries.SearchProfiles:
		return SearchProfiles(ctx, q, p.store.Profiles())
	}

	return nil, cqrs.ErrNoExecutor
}

func (p *pack) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.DeleteMod:
		return DeleteMod(ctx, c, p.store.Mods())
	case *commands.StoreMod:
		return StoreMod(ctx, c, p.store.Mods())
	case *commands.StoreProfile:
		return StoreProfile(ctx, c, p.store.Profiles())
	}

	return cqrs.ErrNoHandler
}

func FromConfig(config *configuration.Config) (Pack, error) {
	storageDriver, err := storageloader.FromConfig(config)
	if err != nil {
		return nil, err
	}

	p := &pack{
		store: storageDriver,
	}

	return p, nil
}
