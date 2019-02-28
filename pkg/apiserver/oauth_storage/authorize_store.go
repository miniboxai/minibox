package oauth_storage

import (
	"errors"

	osin "github.com/RangelReale/osin"
	"minibox.ai/minibox/pkg/models"
)

var authorizeCacheSize = 1024

type AuthorizeStore struct {
	*cacheStore
	db *models.Database
}

func NewAuthorizeStore(db *models.Database) *AuthorizeStore {
	return &AuthorizeStore{
		cacheStore: newCacheStore(authorizeCacheSize),
		db:         db,
	}
}

func cast2authorize(obj interface{}, err error) (*osin.AuthorizeData, error) {
	if err != nil {
		return nil, err
	} else if r, ok := obj.(*osin.AuthorizeData); !ok {
		return nil, errors.New("invalid cast to *osin.AuthorizeData type")
	} else {
		return r, nil
	}
}

func (as *AuthorizeStore) Get(code string) (*osin.AuthorizeData, error) {
	return cast2authorize(as.cacheGet(code, func() (interface{}, error) {
		return nil, osin.ErrNotFound
	}))
}

func (as *AuthorizeStore) Set(code string, data *osin.AuthorizeData) error {
	return as.cacheSet(code, data, func() error {
		return nil
	})
}

func (as *AuthorizeStore) Remove(code string) error {
	return as.cacheRemove(code, func() error {
		return nil
	})
}
