package apiserver

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	osin "github.com/RangelReale/osin"
	"minibox.ai/minibox/pkg/models"
	"minibox.ai/minibox/pkg/sessions"
	"minibox.ai/minibox/pkg/utils"
)

func getParam(r *http.Request, key string) (string, bool) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return "", false
	}

	return keys[0], true
}

// func renderTemplate(w io.Writer, name string, data interface{}) error {
// 	return tmpl.ExecuteTemplate(w, name, data)
// }

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name: name,
	})
}

func encryptCookie(buf []byte) string {
	initOptions()
	crypted, err := utils.TripleDesEncrypt(buf, cryptoKey)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(crypted)
}

func decryptCookie(s string) []byte {
	initOptions()

	buf, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}

	buf, err = utils.TripleDesDecrypt(buf, cryptoKey)
	if err != nil {
		return nil
	}
	return buf
}

func getCurrentUser(session *sessions.Session) *models.User {
	// var usr models.User
	// data, ok := session.Values["cur_user"].([]byte)
	// if !ok {
	// 	return nil
	// }
	// err := json.Unmarshal(data, &usr)
	// if err != nil {
	// 	return nil
	// }
	usr, ok := session.Values["cur_user"].(*models.User)
	if !ok {
		return nil
	}
	return usr
}

func setCurrentUser(session *sessions.Session, usr *models.User) error {
	// buf, err := json.Marshal(usr)
	// if err != nil {
	// 	return err
	// }
	// session.Values["cur_user"] = buf
	session.Values["cur_user"] = usr
	return nil
}

func logout(session *sessions.Session) {
	delete(session.Values, "cur_user")
}

func addBuiltinTemplate(tmpl *template.Template, builtin *template.Template) {
	filename := builtin.Name()
	ext := filepath.Ext(filename)
	if _, err := tmpl.AddParseTree(filename[:len(filename)-len(ext)], builtin.Tree); err != nil {
		fmt.Printf("add parse tree %s\n", err)
	}
}

func getUserByAccessToken(ad *osin.AccessData) (*models.User, error) {
	var (
		userData map[string]interface{}
		ok       bool
		usr      *models.User
	)

	// log.Printf("userData %#v", ad.UserData)

	if ad.UserData == nil {
		ad.UserData = make(map[string]interface{})
	}

	if userData, ok = ad.UserData.(map[string]interface{}); ok {
		if usr, ok = userData["user"].(*models.User); ok {
			return usr, nil
		} else {
			return nil, errors.New("invalid user key in UserData")
		}
	}
	// ad.UserData
	return nil, errors.New("invalid UserData type, must `map[string]interface{}`")
}

func getUserFromModel(ad *osin.AccessData) (*models.User, error) {
	var (
		token models.OAuth2Token
		usr   models.User
	)
	if db.First(&token, "token = ?", ad.AccessToken).RecordNotFound() {
		return nil, errors.New("can load access token")
	}
	if err := db.First(&usr, "id = ?", token.UserID).Error; err != nil {
		return nil, err
	}

	return &usr, nil
}
