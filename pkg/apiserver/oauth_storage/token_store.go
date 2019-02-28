package oauth_storage

import (
	"errors"

	osin "github.com/RangelReale/osin"
	"minibox.ai/pkg/models"
)

var tokenCacheSize = 1024
var ErrNeedSetRefreshStore = errors.New("not set rs *RefreshStore, just call ts.SetRefreshStore(rs)")

type TokenStore struct {
	*cacheStore
	db *models.Database
	rs *RefreshStore
}

func cast2token(obj interface{}, err error) (*osin.AccessData, error) {
	if err != nil {
		return nil, err
	} else if r, ok := obj.(*osin.AccessData); !ok {
		return nil, errors.New("invalid cast to osin.AccessData type")
	} else {
		return r, nil
	}
}

func NewTokenStore(db *models.Database) *TokenStore {
	return &TokenStore{
		cacheStore: newCacheStore(tokenCacheSize),
		db:         db,
	}
}

func (ts *TokenStore) SetRefreshStore(rs *RefreshStore) {
	ts.rs = rs
}

func (ts *TokenStore) Get(code string) (*osin.AccessData, error) {
	return cast2token(ts.cacheGet(code, func() (interface{}, error) {
		var token models.OAuth2Token

		if ts.db.First(&token, "token = ?", code).RecordNotFound() {
			return nil, osin.ErrNotFound
		}

		return token2osin(&token), nil
	}))
}

func (ts *TokenStore) Set(token string, data *osin.AccessData) error {
	err := ts.cacheSet(token, data, func() error {
		tk := osin2token(data)
		return ts.db.Create(tk).Error
	})

	// if data.RefreshToken != "" && ts.rs != nil {
	// 	// s.refresh[data.RefreshToken] = data.AccessToken
	// 	ts.rs.Set(data.RefreshToken, data)
	// }

	return err
}

func (ts *TokenStore) Remove(token string) error {
	return ts.cacheRemove(token, func() error {
		return ts.db.Delete(models.OAuth2Token{}, "token = ?", token).Error
	})
}
