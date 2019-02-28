// Package apiserver API 服务器 组织管理整个业务逻辑，启动 服务器提供核心业务，
// 与 k8s 的沟通， 并提供部署业务， 数据管理服务，以及配置服务等
package apiserver // import "minibox.ai/minibox/pkg/apiserver"

import (
	"log"
	"net/http"
	"os"

	osin "github.com/RangelReale/osin"
	ulogger "github.com/go-http-utils/logger"
	oauth_storage "minibox.ai/minibox/pkg/apiserver/oauth_storage"
	"minibox.ai/minibox/pkg/models"
	oauth2provider "minibox.ai/minibox/pkg/oauth2/provider"
	"minibox.ai/minibox/pkg/sessions"
	"minibox.ai/minibox/pkg/utils/tmpl"
)

// ApiServer 采用了 http.ServeMux 模块载入的方案
type ApiServer struct {
	*http.ServeMux
	callbacks map[string][]CallbackHandler
}

type CallbackHandler func(*ApiServer)

// NewApiServer 创建 ApiServer 实体，稍后用 http.ListenAndServe 启动
func NewApiServer() *ApiServer {
	server := &ApiServer{
		ServeMux: http.NewServeMux(),
	}
	server.init()
	return server
}

type routeHandler func(*ApiServer) error

var db *models.Database

func (svr *ApiServer) init() {
	svr.callbacks = make(map[string][]CallbackHandler)

	svr.staticHandler()
	svr.Handle("/signout", sessions.MiddlewareFunc(func(w http.ResponseWriter, r *http.Request) {
		session := sessions.GetSession(r)
		logout(session)
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}))

	tmpl.LoadTemplates()
	tmpl.AddBuiltinTemplate(buildNavTemplate())

	// db, _ = database.Open("sqlite3", "db/minibox.db")
	// db.LogMode(true)
	err := models.LoadConfig()
	if err != nil {
		log.Fatalf("load database config error: %s", err)
	}
	models.LogMode(true)
	db = models.GetDB()

	Registry(svr)
	EachInitializer(svr)
	models.AutoMigrations(db)
	storage = oauth_storage.NewStorage(db)
}

func (svr *ApiServer) staticHandler() {
	fs := http.FileServer(http.Dir("./static"))
	svr.Handle("/static/", http.StripPrefix("/static/", fs))
	svr.Handle("/website/", http.StripPrefix("/website/", http.FileServer(http.Dir("../website/build/minibox"))))
}

func (svr *ApiServer) beforeRouter(handle CallbackHandler) {
	var (
		callbacks []CallbackHandler
		ok        bool
	)

	if callbacks, ok = svr.callbacks["beforeRouter"]; !ok {
		callbacks = make([]CallbackHandler, 0)
	}

	svr.callbacks["beforeRouter"] = append(callbacks, handle)
}

func (svr *ApiServer) afterRouter(handle CallbackHandler) {
	var (
		callbacks []CallbackHandler
		ok        bool
	)

	if callbacks, ok = svr.callbacks["afterRouter"]; !ok {
		callbacks = make([]CallbackHandler, 0)
	}

	svr.callbacks["afterRouter"] = append(callbacks, handle)
}

func (svr *ApiServer) runCallback(name string) {
	var (
		callbacks []CallbackHandler
		ok        bool
	)

	if callbacks, ok = svr.callbacks[name]; ok {
		for _, callback := range callbacks {
			callback(svr)
		}
	}
}

var storage osin.Storage

// Listen 启动 ApiServer
func Listen(server *ApiServer) error {
	// server.init()
	ns := NewNamespace()
	ds := NewDatasets()

	server.runCallback("beforeRouter")

	server.Handle("/oauth/", http.StripPrefix("/oauth", sessions.Middleware(oauth2provider.NewOAuthProvider(storage))))
	server.Handle("/signin/google/", http.StripPrefix("/signin/google", NewGoogleSignin()))
	server.Handle("/datasets/", http.StripPrefix("/datasets", ds))
	server.Handle("/", sessions.Middleware(homeMiddleware(ns.Entry())))

	server.runCallback("afterRouter")
	log.Println("minibox-apiserver listen on :14000 port")

	return http.ListenAndServe(":14000", ulogger.Handler(server, os.Stdout, ulogger.DevLoggerType))
}
