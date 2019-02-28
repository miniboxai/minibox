package oauth_storage

import (
	"encoding/json"

	"minibox.ai/pkg/models"
)

type (
	Client models.OAuth2Client
	// Authorize models.OAuth2Authorize
	// Token models.OAuth2Token
)

func (cli Client) GetId() string {
	c := (models.OAuth2Client)(cli)
	return c.ID
}

func (cli Client) GetSecret() string {
	c := (models.OAuth2Client)(cli)
	return c.Secret
}

func (cli Client) GetRedirectUri() string {
	c := (models.OAuth2Client)(cli)
	return c.RedirectUri
}

func (cli Client) GetUserData() interface{} {
	var data = make(map[string]interface{})
	c := (models.OAuth2Client)(cli)
	if err := json.Unmarshal(c.UserData, data); err != nil {
		return nil
	}
	return data
}
