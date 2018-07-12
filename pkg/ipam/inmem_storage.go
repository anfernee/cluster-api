package ipam

import "sync"

// inMemStorage is an in memory implementation of interface Storage
type inMemStorage struct {
	// mapping of cluster->ip->address
	free, used map[string]map[string]Address

	sync.Mutex
}

func NewInMemStorage() Storage {
	return &inMemStorage{
		free: make(map[string]map[string]Address),
		used: make(map[string]map[string]Address),
	}
}

func (s *inMemStorage) Add(key string, addresses []Address) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.free[key]; !ok {
		s.free[key] = make(map[string]Address)
		s.used[key] = make(map[string]Address)
	}

	free, used := s.free[key], s.used[key]
	for _, address := range addresses {
		akey := address.Key()
		_, ok1 := free[akey]
		_, ok2 := used[akey]
		if !ok1 && !ok2 {
			free[akey] = address
		}
	}

	return nil
}

func (s *inMemStorage) Remove(key string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.free, key)
	delete(s.used, key)

	return nil
}

func (s *inMemStorage) Allocate(key string) (*Address, error) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.free[key]; !ok {
		return nil, ErrKeyNotFound
	}

	free := s.free[key]
	used := s.used[key]

	if len(free) == 0 {
		return nil, ErrNotEnoughAddress
	}

	var address *Address
	for ip, a := range free {
		used[ip] = a
		delete(free, ip)
		address = &a
		break
	}

	return address, nil
}

func (s *inMemStorage) Release(key, addressKey string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.used[key]; !ok {
		return ErrKeyNotFound
	}

	free := s.free[key]
	used := s.used[key]

	if _, ok := used[addressKey]; !ok {
		return ErrReleaseUnallocatedIP
	}

	free[addressKey] = used[addressKey]
	delete(used, addressKey)

	return nil
}
