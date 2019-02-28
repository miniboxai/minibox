package provider

import osin "github.com/RangelReale/osin"

func setUserData(ar *osin.AuthorizeRequest, key string, value interface{}) {
	var (
		userData map[string]interface{}
		ok       bool
	)

	if ar.UserData == nil {
		userData = make(map[string]interface{})
		ar.UserData = userData
	} else {
		if userData, ok = ar.UserData.(map[string]interface{}); !ok {
			return
		}
	}

	userData[key] = value
	ar.UserData = userData
}
