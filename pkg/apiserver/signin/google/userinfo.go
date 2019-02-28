package google

import (
	"encoding/json"
	"errors"
	"net/http"

	oauth2 "golang.org/x/oauth2"
)

const GoogleV3UserInfo = "https://www.googleapis.com/oauth2/v3/userinfo"

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func FetchUser(client *http.Client, tok *oauth2.Token) (*UserInfo, error) {
	var user UserInfo

	email, err := client.Get(GoogleV3UserInfo)
	if err != nil {
		return nil, err
	}

	defer email.Body.Close()
	dec := json.NewDecoder(email.Body)
	if dec == nil {
		return nil, errors.New("can't create decoder from userinfo Body stream")
	}

	if err := dec.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
