package actions

import (
	"context"

	"github.com/danielkrainas/shexd/api/v1"
	"github.com/danielkrainas/shexd/commands"
	"github.com/danielkrainas/shexd/queries"
	"github.com/danielkrainas/shexd/storage"
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

func StoreProfile(ctx context.Context, c *commands.StoreProfile, profiles storage.ProfileStore) error {
	p := c.Profile
	return profiles.Store(p, c.New)
}

func SearchProfiles(ctx context.Context, q *queries.SearchProfiles, profiles storage.ProfileStore) ([]*v1.RemoteProfile, error) {
	return profiles.FindMany(&storage.ProfileFilters{})
}

func GetModVersionList(ctx context.Context, q *queries.GetModVersionList, mods storage.ModStore) ([]string, error) {
	return mods.Versions(q.Token)
}
