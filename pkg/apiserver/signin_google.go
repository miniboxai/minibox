package apiserver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/schema"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"minibox.ai/minibox/pkg/apiserver/signin/google"
	"minibox.ai/minibox/pkg/models"
	"minibox.ai/minibox/pkg/sessions"
	"minibox.ai/minibox/pkg/utils/tmpl"
)

// Google 登陆服务模块， 我们采用了多个 ServeMux 组合服务的方案，这让我们可以通过下面语句
//	server.Handle("/signin/google/", http.StripPrefix("/signin/google", NewGoogleSignin()))
// 添加登陆到我们的服务当中去， 这样我们可以组合不同的服务类型
type GoogleSignin struct {
	*http.ServeMux
}

var Prefix = "/signin/google"

const cookieUserInfo = "__callback_userinfo"
const signinReferralUrl = "_signin_referral_url"

type byUser int

const (
	byName byUser = iota
	byEmail
	byNamespace
	byMobile
)

var cryptoKey []byte

func NewGoogleSignin() *GoogleSignin {
	mux := &GoogleSignin{
		ServeMux: http.NewServeMux(),
	}
	mux.init()
	return mux
}

func (sign *GoogleSignin) init() {
	sign.HandleFunc("/", googleSignin)

	sign.Handle("/callback", sessions.MiddlewareFunc(googleCallback))
	sign.HandleFunc("/new_namespace", googleNewNamespace)
	sign.Handle("/create_namespace", sessions.MiddlewareFunc(googleCreateNamespace))
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func googleSignin(w http.ResponseWriter, r *http.Request) {
	scope := "profile"
	if !google.ValidConfig() {
		tmpl.RenderTemplate(w, "google_config.html.tmpl", map[string]interface{}{
			"DOC_URL": "",
		})
		return	
	}

	tmpl.RenderTemplate(w, "login.html.tmpl", map[string]interface{}{
		"LoginURL": google.LoginURL(scope),
	})
}

func googleCallback(w http.ResponseWriter, r *http.Request) {
	// session := sessions.Default(c)
	//    retrievedState := session.Get("state")
	//    if retrievedState != c.Query("state") {
	//        c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
	//        return
	//    }
	var (
		code string
		ok   bool
		err  error
		tok  *oauth2.Token
	)
	r.ParseForm()
	if code, ok = getParam(r, "code"); !ok {
		http.Error(w, "missing code param", http.StatusBadRequest)
		return
	}

	// Handle the exchange code to initiate a transport.
	tok, err = google.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Construct the client.
	client := google.NewClient(oauth2.NoContext, tok)
	user, err := google.FetchUser(client, tok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if validateGoogleUserExists(user) {
		if session := sessions.GetSession(r); session != nil {
			usr, err := loadUserBy(byEmail, user.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			setCurrentUser(session, usr)
			session.Save(r, w)
			// session.Values["cur_user"] = usr
		}

		http.Redirect(w, r, getCallbackUrl(r), http.StatusFound)
		clearCookie(w, signinReferralUrl)

	} else {
		url := path.Join(Prefix, "new_namespace")
		j, _ := json.Marshal(user)

		http.SetCookie(w, &http.Cookie{
			Name:  cookieUserInfo,
			Value: encryptCookie(j),
		})
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func validateGoogleUserExists(googleUser *google.UserInfo) bool {
	var user models.User
	return !db.First(&user, "email = ?", googleUser.Email).RecordNotFound()
}

func loadUserBy(by byUser, value string) (*models.User, error) {
	var user models.User

	switch by {
	case byEmail:
		if err := db.First(&user, "email = ?", value).Error; err != nil {
			return nil, err
		}
	case byNamespace:
		if err := db.First(&user, "namespace = ?", value).Error; err != nil {
			return nil, err
		}
	case byMobile:
		if err := db.First(&user, "mobile = ?", value).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func getCallbackUrl(r *http.Request) string {
	cookie, err := r.Cookie(signinReferralUrl)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func googleNewNamespace(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieUserInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := decryptCookie(cookie.Value)
	if buf == nil {
		http.Error(w, "invalid goolge userinfo", http.StatusInternalServerError)
		return
	}
	var user google.UserInfo
	if err := json.Unmarshal(buf, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var username string
	if pos := strings.Index(user.Email, "@"); pos > 0 {
		username = user.Email[:pos]
	}

	form := &CreateNamespaceForm{
		Action:     path.Join(Prefix, "/create_namespace"),
		Namespace:  username,
		Name:       user.Name,
		Email:      user.Email,
		Sub:        user.Sub,
		Gender:     user.Gender,
		FamilyName: user.FamilyName,
		GivenName:  user.GivenName,
		Profile:    user.Profile,
		Picture:    user.Picture,
	}
	// clearCookie(w, cookieUserInfo)
	tmpl.RenderTemplate(w, "new_namespace.html.tmpl", map[string]interface{}{
		"Form": form,
	})
}

func googleCreateNamespace(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var (
		form CreateNamespaceForm
		user models.User
	)

	decoder := schema.NewDecoder()
	// r.PostForm is a map of our POST form values
	log.Printf("PostForm: %#v", r.PostForm)
	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errs := validateForm(&form); errs != nil {
		for _, err := range rangeValidationErrors(errs) {
			fmt.Fprintf(w, "error: %s", err)
		}

		return
	}

	user.Namespace = form.Namespace
	user.Name = form.Name
	user.Email = form.Email
	user.Avatar = form.Picture
	user.Provider = models.Provider{
		Name:  "google",
		SubID: form.Sub,
	}

	if err = db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session := sessions.GetSession(r); session != nil {
		setCurrentUser(session, &user)
		session.Save(r, w)
	}

	tmpl.RenderTemplate(w, "create_namespace.html.tmpl", &user)
	clearCookie(w, signinReferralUrl)
}

func initOptions() {
	cfg := viper.Sub("server.signin.cookie.crypto")

	if cfg != nil {
		keys := cfg.GetString("key")
		if len(keys) >= 24 {
			cryptoKey = []byte(keys)[:24]
		} else {
			panic("server.signin.cookie.crypto.key must large then 24 bytes")
		}
	} else {
		cryptoKey = []byte("example key 1234567890xx")
	}
}

func init() {
	gob.Register(&models.User{})
}
