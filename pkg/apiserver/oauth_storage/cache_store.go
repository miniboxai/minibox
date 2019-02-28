package oauth_storage

import lru "github.com/hashicorp/golang-lru"

type cacheStore struct {
	cache *lru.Cache
}

type getCallback func() (interface{}, error)
type setCallback func() error
type removeCallback func() error

func newCacheStore(size int) *cacheStore {
	cache, _ := lru.New(size)
	return &cacheStore{
		cache: cache,
	}
}

func (cs *cacheStore) cacheGet(id string, reallyGet getCallback) (interface{}, error) {
	if obj, ok := cs.cache.Get(id); ok { // has cache
		return obj, nil
	} else if obj, err := reallyGet(); err != nil { // callback really data fetch
		return nil, err
	} else { // store to cache again
		cs.cache.Add(id, obj)
		return obj, nil
	}
}

func (cs *cacheStore) cacheSet(id string, val interface{}, reallySet setCallback) error {
	if err := reallySet(); err != nil {
		return err
	} else {
		cs.cache.Add(id, val)
		return nil
	}
}

func (cs *cacheStore) cacheRemove(id string, reallyRemove removeCallback) error {
	if err := reallyRemove(); err != nil {
		return err
	} else {
		cs.cache.Remove(id)
		return nil
	}
}
