package oauth_storage

import (
	"errors"
	"log"

	osin "github.com/RangelReale/osin"
	"minibox.ai/pkg/models"
)

var clientCacheSize = 128

type ClientStore struct {
	*cacheStore
	db *models.Database
}

func NewClientStore(db *models.Database) *ClientStore {
	return &ClientStore{
		cacheStore: newCacheStore(clientCacheSize),
		db:         db,
	}
}

func cast2client(obj interface{}, err error) (osin.Client, error) {
	if err != nil {
		return nil, err
	} else if r, ok := obj.(Client); !ok {
		log.Printf("obj %#v", obj)
		return nil, errors.New("invalid cast to osin.Client type")
	} else {
		return r, nil
	}
}

func (cs *ClientStore) Get(id string) (osin.Client, error) {
	return cast2client(cs.cacheGet(id, func() (interface{}, error) {
		var (
			client models.OAuth2Client
			user   models.User
		)

		if cs.db.Model(&user).First(&client, "id = ?", id).RecordNotFound() {
			return nil, osin.ErrNotFound
		}
		return Client(client), nil
	}))
}

func (cs *ClientStore) Set(id string, client osin.Client) error {
	return cs.cacheSet(id, client, func() error {
		var cli = models.OAuth2Client{
			ID:          client.GetId(),
			Secret:      client.GetSecret(),
			RedirectUri: client.GetRedirectUri(),
		}
		return cs.db.FirstOrCreate(&cli).Error
	})
}

func (cs *ClientStore) Create(id string, client osin.Client) error {
	return cs.cacheSet(id, client, func() error {
		var cli = models.OAuth2Client{
			ID:          client.GetId(),
			Secret:      client.GetSecret(),
			RedirectUri: client.GetRedirectUri(),
		}
		return cs.db.Create(&cli).Error
	})
}
