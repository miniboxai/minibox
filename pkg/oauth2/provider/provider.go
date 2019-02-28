package provider

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"golang.org/x/oauth2"
	"minibox.ai/pkg/apiserver/signin/google"
	oclient "minibox.ai/pkg/oauth2/client"
	"minibox.ai/pkg/sessions"
	"minibox.ai/pkg/utils/tmpl"
)

var Prefix = "/oauth"

type LoginHandler func(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool

var defaultLogin = HandleLoginPage

type OAuthProvider struct {
	*http.ServeMux
	store    osin.Storage
	jwtToken *AccessTokenGenJWT
	server   *osin.Server
	login    LoginHandler
}

func NewOAuthProvider(store osin.Storage) *OAuthProvider {
	mux := &OAuthProvider{
		ServeMux: http.NewServeMux(),
		store:    store,
		login:    defaultLogin,
	}
	mux.init()
	return mux
}

func (mux *OAuthProvider) SetStorage(store osin.Storage) {
	mux.store = store
}

func HandleLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()

	// if r.Method == "POST" {
	session := sessions.GetSession(r)
	cur_user := sessions.GetCurrentUser(session)
	if cur_user != nil { // 有用户信息
		return true

	}
	scope := "profile"

	if !google.ValidConfig() {
		tmpl.RenderTemplate(w, "google_config.html.tmpl", map[string]interface{}{
			"DOC_URL": "https://miniboxai.github.io/website/docs/en/doc1.html",
		})
		return false
	}	

	tmpl.RenderTemplate(w, "login.html.tmpl", map[string]interface{}{
		"LoginURL": google.LoginURL(scope),
		// "LoginURL": fmt.Sprintf("/oauth/authorize?%s", r.URL.RawQuery),
	})

	return false
}

func HandleLoginPage2(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {

	if r.Method == "POST" && r.FormValue("login") == "test" && r.FormValue("password") == "test" {
		return true
	}
	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/oauth/authorize?%s\" method=\"POST\">", r.URL.RawQuery)))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return false
}

func HandleGrantAccess(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()

	if ar.Authorized {
		session := sessions.GetSession(r)
		cur_user := sessions.GetCurrentUser(session)
		setUserData(ar, "user", cur_user)

		if r.FormValue("dogrant") == "1" {
			return true
		} else {
			tmpl.RenderTemplate(w, "grant_access.html.tmpl", cur_user)
		}
	}
	return false
}

func (mux *OAuthProvider) init() {
	config := osin.NewServerConfig()
	// goauth2 checks errors using status codes
	config.ErrorStatusCode = 401

	server := osin.NewServer(config, mux.store)
	mux.server = server
	tokenGen, _ := loadTokenGenJWT()
	tokenGen.SetClientCallback(func(cid string, data *osin.AccessData) (err error) {
		data.Client, err = mux.store.GetClient(cid)
		return err
	})

	mux.jwtToken = tokenGen

	client := oclient.NewClient("1234", "aabbccdd", "http://localhost:14000/oauth/appauth/code")

	// Authorization code endpoint
	mux.Handle("/authorize", sessions.Middleware(mux.authorizeHandle()))

	// Access token endpoint
	mux.Handle("/token", sessions.Middleware(mux.tokenHandle()))

	mux.HandleFunc("/jwt", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()
		mux.downloadJwt(resp, r)

		osin.OutputJSON(resp, w, r)
	})

	// Information endpoint
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ir := server.HandleInfoRequest(resp, r); ir != nil {
			server.FinishInfoRequest(resp, r, ir)
		}
		osin.OutputJSON(resp, w, r)
	})

	// scope := oauth2.SetAuthURLParam("scope", "cli")
	// Application home endpoint
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte(`<html><body>`))
		// w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Login</a><br/>", client.AuthCodeURL(""))))
		// w.Write([]byte("</body></html>"))

		tmpl.RenderTemplate(w, "login.html.tmpl", map[string]interface{}{
			"LoginURL": client.AuthCodeURL(""),
		})
	})

	// Application destination - CODE
	mux.HandleFunc("/appauth/code", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		code := r.FormValue("code")

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CODE<br/>"))
		defer w.Write([]byte("</body></html>"))

		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		var jr *oauth2.Token
		var err error

		// if parse, download and parse json
		if r.FormValue("doparse") == "1" {
			jr, err = client.Exchange(oauth2.NoContext, code)
			if err != nil {
				jr = nil
				w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", err)))
			}
		}

		// show json access token
		if jr != nil {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", jr.AccessToken)))
			if jr.RefreshToken != "" {
				w.Write([]byte(fmt.Sprintf("REFRESH TOKEN: %s<br/>\n", jr.RefreshToken)))
			}
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		cururl := *r.URL
		curq := cururl.Query()
		curq.Add("doparse", "1")
		cururl.RawQuery = curq.Encode()
		w.Write([]byte(fmt.Sprintf("<a href=\"%s%s\">Download Token</a><br/>", Prefix, cururl.String())))
	})
}

func (mux *OAuthProvider) authorizeHandle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := mux.server.NewResponse()
		defer resp.Close()

		if ar := mux.server.HandleAuthorizeRequest(resp, r); ar != nil {
			if !mux.login(ar, w, r) {
				return
			}

			ar.Authorized = true
			if !HandleGrantAccess(ar, w, r) {
				return
			}

			mux.server.FinishAuthorizeRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	})
}

func (mux *OAuthProvider) tokenHandle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := mux.server.NewResponse()
		defer resp.Close()
		ar := mux.server.HandleAccessRequest(resp, r)

		if ar != nil {
			ar.Authorized = true
			if ar.Scope == "cli" {
				ar.Expiration = 31556926
			}

			mux.server.FinishAccessRequest(resp, r, ar)
		}

		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}

		osin.OutputJSON(resp, w, r)
	})
}

func (mux *OAuthProvider) downloadJwt(resp *osin.Response, r *http.Request) error {
	r.ParseForm()
	bearer := osin.CheckBearerAuth(r)
	if bearer == nil {
		// s.setErrorAndLog(w, osin.E_INVALID_REQUEST, nil, "handle_jwt_request=%s", "bearer is nil")
		// return osin.E_INVALID_REQUEST
		// resp.InternalError = osin.E_INVALID_REQUEST
		resp.SetError(osin.E_INVALID_REQUEST, "")

		// resp.SetError("handle_jwt_request=bearer is nil", "")
		return nil
	}

	code := bearer.Code

	if code == "" {
		resp.SetError(osin.E_INVALID_REQUEST, "")
		return nil
	}

	var err error

	// load access data
	ad, err := mux.store.LoadAccess(code)
	if err != nil {
		// resp.Set
		// resp.InternalError = osin.E_INVALID_REQUEST
		resp.SetError(osin.E_INVALID_REQUEST, "")
		// s.setErrorAndLog(w, E_INVALID_REQUEST, err, "handle_info_request=%s", "failed to load access data")
		return nil
	}

	estr, err := mux.jwtToken.EncryptoMap(map[string]interface{}{
		"client_id":     ad.Client.GetId(),
		"access_token":  ad.AccessToken,
		"refresh_token": ad.RefreshToken,
		"expires_in":    ad.ExpireAt().Unix(),
	})
	resp.Output["jwt"] = estr
	return nil
}

func (mux *OAuthProvider) SetDefaultLogin(handle LoginHandler) {
	mux.login = handle
}

func mergeOutput(a, b map[string]interface{}) map[string]interface{} {
	for key, value := range b {
		a[key] = value
	}

	return a
}

func jwtToken(tokenGen *AccessTokenGenJWT, mapClaim osin.ResponseData, client osin.Client) error {
	token, _ := mapClaim["access_token"].(string)
	refreshToken, _ := mapClaim["refresh_token"].(string)
	expires_in, _ := mapClaim["expires_in"].(string)
	delete(mapClaim, "access_token")
	delete(mapClaim, "refresh_token")
	delete(mapClaim, "expires_in")
	delete(mapClaim, "token_type")

	estr, err := tokenGen.EncryptoMap(map[string]interface{}{
		"client_id":     client.GetId(),
		"access_token":  token,
		"refresh_token": refreshToken,
		"expires_in":    expires_in,
	})
	if err != nil {
		return err
	}
	mapClaim["jwt"] = estr
	return nil
}
