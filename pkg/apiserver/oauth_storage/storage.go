package oauth_storage

import (
	"fmt"

	osin "github.com/RangelReale/osin"
	"minibox.ai/pkg/models"
)

type Storage struct {
	clients     *ClientStore
	authorize   *AuthorizeStore
	access      *TokenStore
	refresh     *RefreshStore
	initializes []StorageHandler
}

type StorageHandler func(*Storage)

func NewStorage(db *models.Database) *Storage {

	r := &Storage{
		initializes: make([]StorageHandler, 0),
	}

	r.clients = NewClientStore(db)
	r.authorize = NewAuthorizeStore(db)
	r.access = NewTokenStore(db)
	r.refresh = NewRefreshStore(db)
	// r.access.SetRefreshStore(r.refresh)
	// r.refresh.SetTokenStore(r.access)

	r.clients.Set("1234", Client(models.OAuth2Client{
		ID:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/oauth/appauth",
	}))

	r.clients.Set("12345", Client(models.OAuth2Client{
		ID:          "12345",
		Secret:      "asdfasdf",
		RedirectUri: "http://localhost:14001/auth",
	}))

	r.init()
	return r
}

func (s *Storage) init() {
	EachInitializer(s)
}

func (s *Storage) Clone() osin.Storage {
	return s
}

func (s *Storage) Close() {
}

func (s *Storage) GetClient(id string) (osin.Client, error) {
	fmt.Printf("GetClient: %s\n", id)
	return s.clients.Get(id)
}

func (s *Storage) SetClient(id string, client osin.Client) error {
	fmt.Printf("SetClient: %s\n", id)
	return s.clients.Set(id, client)
}

func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) error {
	fmt.Printf("SaveAuthorize: %s\n", data.Code)
	return s.authorize.Set(data.Code, data)
}

func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	fmt.Printf("LoadAuthorize: %s\n", code)
	return s.authorize.Get(code)
}

func (s *Storage) RemoveAuthorize(code string) error {
	fmt.Printf("RemoveAuthorize: %s\n", code)
	return s.authorize.Remove(code)
}

func (s *Storage) SaveAccess(data *osin.AccessData) error {
	fmt.Printf("SaveAccess: %s\n", data.AccessToken)

	if err := s.access.Set(data.AccessToken, data); err != nil {
		return err
	}

	if data.RefreshToken != "" {
		s.refresh.Set(data.RefreshToken, data)
	}
	return nil
}

func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadAccess: %s\n", code)
	return s.access.Get(code)
}

func (s *Storage) RemoveAccess(code string) error {
	fmt.Printf("RemoveAccess: %s\n", code)
	return s.access.Remove(code)
}

func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadRefresh: %s\n", code)
	var (
		token *osin.AccessData
		err   error
	)

	if token, err = s.refresh.Get(code); err != nil {
		return nil, err
	}

	return s.access.Get(token.AccessToken)
}

func (s *Storage) RemoveRefresh(code string) error {
	fmt.Printf("RemoveRefresh: %s\n", code)
	return s.refresh.Remove(code)
}
