package client

import (
	"errors"

	"github.com/spf13/viper"
	oauth2 "golang.org/x/oauth2"
	"minibox.ai/minibox/pkg/database"
	"minibox.ai/minibox/pkg/models"
)

// config :=  client.From("1234", "secret")
var ErrNotFoundClient = errors.New("can't found available Client")
var defaultEndpoints *oauth2.Endpoint

type Config struct {
	db *database.Database
}

func New(db *database.Database) *Config {
	return &Config{
		db: db,
	}
}

func NewClient(id, secret, redirectUri string) *oauth2.Config {
	initOptions()

	return &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     *defaultEndpoints,
		RedirectURL:  redirectUri,
	}
}

func (cfg *Config) LoadClient(id, secret string) (*oauth2.Config, error) {
	initOptions()
	var client models.OAuth2Client

	if cfg.db.First(&client, "id = ? and secret = ?", id, secret).RecordNotFound() {
		return nil, ErrNotFoundClient
	}

	return &oauth2.Config{
		ClientID:     client.ID,
		ClientSecret: client.Secret,
		Endpoint:     *defaultEndpoints,
		RedirectURL:  client.RedirectUri,
	}, nil
}

func initOptions() {
	if defaultEndpoints != nil {
		return
	}
	var (
		authURL  = "http://localhost:14000/oauth/authorize"
		tokenURL = "http://localhost:14000/oauth/token"
		cfg      = viper.Sub("server.oauth2.endpoint")
	)

	if cfg != nil {
		cfg.SetDefault("auth_url", authURL)
		cfg.SetDefault("token_url", tokenURL)
		authURL = cfg.GetString("auth_url")
		tokenURL = cfg.GetString("token_url")
	} else {
		cfg = viper.New()
	}
	defaultEndpoints = &oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}
}

// func (cfg *Config) NewClientDefault(id, secret string) (*oauth2.Config, error) {
// 	var client = models.OAuth2Client{
// 		ID:     id,
// 		Secret: secret,
// 		Endpoint: models.OAuth2Endpoint{
// 			AuthURL:  "",
// 			TokenURL: "",
// 		},
// 	}

// 	return &client, nil
// }
