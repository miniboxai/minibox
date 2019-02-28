package oauth_storage

import (
	"errors"

	osin "github.com/RangelReale/osin"
	"minibox.ai/pkg/models"
)

var refreshCacheSize = 1024

type RefreshStore struct {
	*cacheStore
	db *models.Database
	ts *TokenStore
}

var ErrNeedSetTokenStore = errors.New("not set ts *TokenStore, just call rs.SetTokenStore(ts)")

func NewRefreshStore(db *models.Database) *RefreshStore {
	return &RefreshStore{
		cacheStore: newCacheStore(refreshCacheSize),
		db:         db,
	}
}

func (rs *RefreshStore) SetTokenStore(ts *TokenStore) {
	rs.ts = ts
}

func cast2string(obj interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	} else if s, ok := obj.(string); !ok {
		return "", errors.New("invalid cast to string type")
	} else {
		return s, nil
	}
}

func (rs *RefreshStore) Get(code string) (*osin.AccessData, error) {
	var token models.OAuth2Token

	return cast2token(rs.cacheGet(code, func() (interface{}, error) {

		if rs.db.First(&token, "refresh = ?", code).RecordNotFound() {
			return "", osin.ErrNotFound
		}

		// return token.Token, nil
		return token2osin(&token), nil
	}))

	// if rs.ts != nil { // 如果有 token 会同时载入 AccessToken 数据
	// 	return cast2token(rs.ts.cacheGet(tcode, func() (interface{}, error) {
	// 		return token, error
	// 	}))
	// } else {
	// 	return nil, ErrNeedSetTokenStore
	// }
}

func (rs *RefreshStore) Set(code string, data *osin.AccessData) error {
	return rs.cacheSet(code, data.AccessToken, func() error {
		return nil
	})
}

func (rs *RefreshStore) Remove(code string) error {
	return rs.cacheRemove(code, func() error {
		return nil
	})
}
