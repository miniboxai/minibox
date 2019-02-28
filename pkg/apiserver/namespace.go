package apiserver

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
)

type Namespace struct {
	*http.ServeMux
}

var nsKey = struct{}{}

func NewNamespace() *Namespace {
	mux := &Namespace{
		ServeMux: http.NewServeMux(),
	}
	mux.init()
	return mux
}

func (ns *Namespace) init() {

}

func (ns *Namespace) IndexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if name, ok := r.Context().Value(nsKey).(string); ok {
			fmt.Fprintf(w, "hello %s\n", name)
		}
	})
}

func (ns *Namespace) Entry() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			name  string
			index bool
		)
		dir := html.EscapeString(r.URL.Path)
		if pos := strings.Index(dir[1:], "/"); pos > 4 {
			name = dir[1:pos]
		} else {
			name = dir[1:]
			index = true
		}

		_, err := loadUserBy(byNamespace, name)
		if err != nil {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		r = setNamespace(r, name)
		if index {
			ns.IndexHandler().ServeHTTP(w, r)
		} else {
			ns.ServeHTTP(w, r)
		}
	})
}

func setNamespace(r *http.Request, name string) *http.Request {
	ctx := context.WithValue(r.Context(), nsKey, name)
	return r.WithContext(ctx)
}
