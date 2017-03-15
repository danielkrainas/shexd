package inmemory

import (
	"sync"

	"github.com/danielkrainas/shexd/api/v1"
	"github.com/danielkrainas/shexd/storage"
)

type profileStore struct {
	m        sync.Mutex
	profiles []*v1.RemoteProfile
}

func (s *profileStore) Store(u *v1.RemoteProfile, isNew bool) error {
	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if isNew {
		for i, u2 := range s.profiles {
			if u2.Name == u.Name {
				s.profiles[i] = u
				found = true
				break
			}
		}
	}

	if !found {
		s.profiles = append(s.profiles, u)
	}

	return nil
}

/*func (s *profileStore) Delete(token *v1.NameVersionToken) error {
	s.m.Lock()
	defer s.m.Unlock()
	for i, u := range s.profiles {
		if u.Name == token.Name && fmt.Sprint(u.Version) == token.Version {
			s.profiles = append(s.profiles[:i], s.profiles[i+1:]...)
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *profileStore) Find(token *v1.NameVersionToken) (*v1.RemoteProfile, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.profiles {
		if u.Name == token.Name && fmt.Sprint(u.Version) == token.Version {
			return u, nil
		}
	}

	return nil, nil
}

func (s *profileStore) Count(f *storage.ProfileFilters) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.profiles), nil
}*/

func (s *profileStore) FindMany(f *storage.ProfileFilters) ([]*v1.RemoteProfile, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.profiles[:], nil
}
