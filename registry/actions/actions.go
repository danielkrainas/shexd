package actions

import (
	"context"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/registry/commands"
	"github.com/danielkrainas/shex/registry/queries"
	"github.com/danielkrainas/shex/registry/storage"
)

func DeleteMod(ctx context.Context, c *commands.DeleteMod, mods storage.ModStore) error {
	return mods.Delete(c.Token)
}

func StoreMod(ctx context.Context, c *commands.StoreMod, mods storage.ModStore) error {
	u := c.Mod
	return mods.Store(u, c.New)
}

func FindMod(ctx context.Context, q *queries.FindMod, mods storage.ModStore) (*v1.ModInfo, error) {
	return mods.Find(q.Token)
}

func CountMods(ctx context.Context, q *queries.CountMods, mods storage.ModStore) (int, error) {
	return mods.Count(&storage.ModFilters{})
}

func SearchMods(ctx context.Context, q *queries.SearchMods, mods storage.ModStore) ([]*v1.ModInfo, error) {
	return mods.FindMany(&storage.ModFilters{})
}
