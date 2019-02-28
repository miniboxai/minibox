package google

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

func initConfig() {

	if conf != nil {
		return
	}

	var cid, csecret, redirectUrl string
	cfg := viper.Sub("server.oauth2.signin.google")
	log.Printf("cfg: %#v", cfg)
	if cfg != nil {
		cfg.SetDefault("RedirectURL", "http://localhost:14000/signin/google/callback")
		cid = cfg.GetString("ClientID")
		csecret = cfg.GetString("Secret")
		redirectUrl = cfg.GetString("RedirectURL")
	}

	conf = &oauth2.Config{
		ClientID:     cid,
		ClientSecret: csecret,
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func ValidConfig() bool {
	cfg := viper.Sub("server.oauth2.signin.google")
	return cfg != nil
}

func LoginURL(state string) string {
	initConfig()
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return conf.AuthCodeURL(state)
}

func Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	initConfig()

	return conf.Exchange(ctx, code, opts...)
}

func NewClient(ctx context.Context, src *oauth2.Token) *http.Client {
	initConfig()

	return conf.Client(ctx, src)
}
