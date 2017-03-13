package inmemory

import (
	"fmt"
	"sync"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/registry/storage"
)

type modStore struct {
	m    sync.Mutex
	mods []*v1.ModInfo
}

func (s *modStore) Store(u *v1.ModInfo, isNew bool) error {
	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if isNew {
		for i, u2 := range s.mods {
			if u2.Name == u.Name && u2.Version == u.Version {
				s.mods[i] = u
				found = true
				break
			}
		}
	}

	if !found {
		s.mods = append(s.mods, u)
	}

	return nil
}

func (s *modStore) Delete(token *v1.NameVersionToken) error {
	s.m.Lock()
	defer s.m.Unlock()
	for i, u := range s.mods {
		if u.Name == token.Name && fmt.Sprint(u.Version) == token.Version {
			s.mods = append(s.mods[:i], s.mods[i+1:]...)
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *modStore) Find(token *v1.NameVersionToken) (*v1.ModInfo, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.mods {
		if u.Name == token.Name && fmt.Sprint(u.Version) == token.Version {
			return u, nil
		}
	}

	return nil, nil
}

func (s *modStore) Count(f *storage.ModFilters) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.mods), nil
}

func (s *modStore) FindMany(f *storage.ModFilters) ([]*v1.ModInfo, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.mods[:], nil
}
