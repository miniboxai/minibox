package sessions

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

var (
	storeCookieKey = "_minibox"
	store          sessions.Store
	defaultSession = "minibox"
	sessionKey     = struct{}{}
)

type Session struct {
	*sessions.Session
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// 处理 Session 自动保存的中间件
func MiddlewareFunc(next HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initOptions()

		s, _ := store.Get(r, defaultSession)
		session := &Session{s}

		ctx := context.WithValue(r.Context(), sessionKey, session)
		r = r.WithContext(ctx)
		next(w, r)
		// session.Save(r, w)
	})
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initOptions()

		s, _ := store.Get(r, defaultSession)
		session := &Session{s}
		ctx := context.WithValue(r.Context(), sessionKey, session)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		// session.Save(r, w)
	})
}

func GetSession(r *http.Request) *Session {
	return r.Context().Value(sessionKey).(*Session)
}

func initOptions() {
	if store != nil {
		return
	}

	cfg := viper.Sub("server.sessions")
	if cfg != nil {
		cfg.SetDefault("cookieKey", storeCookieKey)
		storeCookieKey = cfg.GetString("cookieKey")
	}

	store = sessions.NewCookieStore([]byte(storeCookieKey))
}

func init() {
	gob.Register(&Session{})
}
