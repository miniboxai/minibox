package oauth_storage

import (
	"encoding/json"
	"fmt"
	"log"

	osin "github.com/RangelReale/osin"
	"minibox.ai/pkg/models"
)

func oauth2osin(oauth *models.OAuth2Authorize) *osin.AuthorizeData {
	var data osin.AuthorizeData

	data.Code = oauth.Code
	data.ExpiresIn = oauth.ExpiresIn
	data.Scope = oauth.Scope
	data.RedirectUri = oauth.RedirectUri
	data.State = oauth.State
	data.CreatedAt = oauth.CreatedAt
	data.CodeChallenge = oauth.CodeChallenge
	json.Unmarshal(oauth.UserData, &data.UserData)
	data.CodeChallengeMethod = oauth.CodeChallengeMethod
	return &data
}

func osin2auuth(authorize *osin.AuthorizeData) *models.OAuth2Authorize {
	var author models.OAuth2Authorize

	author.Code = authorize.Code
	author.ExpiresIn = authorize.ExpiresIn
	author.Scope = authorize.Scope
	author.RedirectUri = authorize.RedirectUri
	author.State = authorize.State
	author.CreatedAt = authorize.CreatedAt
	author.CodeChallenge = authorize.CodeChallenge
	author.UserData, _ = json.Marshal(authorize.UserData)
	author.CodeChallengeMethod = authorize.CodeChallengeMethod
	return &author
}

func token2osin(token *models.OAuth2Token) *osin.AccessData {
	var (
		data     osin.AccessData
		userData = make(map[string]interface{})
		err      error
	)
	// Client information
	data.Client = Client(token.Client)
	data.AccessToken = token.Token
	data.RefreshToken = token.Refresh
	data.ExpiresIn = token.ExpiresIn
	data.Scope = token.Scope
	data.RedirectUri = token.RedirectUri
	data.CreatedAt = token.CreatedAt

	if err = json.Unmarshal(token.UserData, &userData); err == nil {
		if usr, ok := userData2User(userData); ok {
			userData["user"] = usr
			data.UserData = userData
		}
	} else {
		log.Printf("Unmarshal error %s", err)
	}

	return &data
}

func osin2token(data *osin.AccessData) *models.OAuth2Token {
	var (
		token    models.OAuth2Token
		userData = make(map[string]interface{})
		ok       bool
		err      error
	)

	// token.Client =
	token.Token = data.AccessToken
	token.Refresh = data.RefreshToken
	token.ExpiresIn = data.ExpiresIn
	token.Scope = data.Scope
	token.RedirectUri = data.RedirectUri
	token.CreatedAt = data.CreatedAt
	token.ClientID = data.Client.GetId()

	if userData, ok = data.UserData.(map[string]interface{}); ok {
		guessUser(&token, userData)

		if token.UserData, err = json.Marshal(&userData); err != nil {
			log.Printf("marshal data.UserData error %s", err)
		}
	}

	// token.Client = data.Client

	return &token
}

func guessUser(token *models.OAuth2Token, userData map[string]interface{}) bool {
	if usr, ok := userData["user"].(*models.User); ok {
		token.UserID = int(usr.ID)
		return true
	}
	return false
}

func userData2User(userData map[string]interface{}) (*models.User, bool) {
	var (
		umap = make(map[string]interface{})
		usr  models.User
		ok   bool
		b    []byte
		err  error
	)
	if umap, ok = userData["user"].(map[string]interface{}); !ok {
		return nil, false
	}

	if b, err = json.Marshal(umap); err != nil {
		return nil, false
	}

	if err = json.Unmarshal(b, &usr); err != nil {
		return nil, false
	}
	return &usr, true
}

type SafeMap map[interface{}]interface{}

func (m SafeMap) UnmarshalJSON(b []byte) error {
	var tmap = make(map[string]interface{})
	if err := json.Unmarshal(b, tmap); err != nil {
		return err
	}
	for key, val := range tmap {
		m[key] = val
	}
	return nil
}

func (m SafeMap) MarshalJSON() ([]byte, error) {
	var tmap = make(map[string]interface{})
	for key, val := range m {

		tmap[fmt.Sprintf("%s", key)] = val
	}
	return json.Marshal(tmap)
}
