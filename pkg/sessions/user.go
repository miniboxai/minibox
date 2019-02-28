package sessions

import (
	"minibox.ai/pkg/models"
)

func GetCurrentUser(session *Session) *models.User {
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

func SetCurrentUser(session *Session, usr *models.User) error {
	// buf, err := json.Marshal(usr)
	// if err != nil {
	// 	return err
	// }
	// session.Values["cur_user"] = buf
	session.Values["cur_user"] = usr
	return nil
}
