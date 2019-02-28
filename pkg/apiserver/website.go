package apiserver

import "net/http"

func mountWebsiteHandler(hostdir string) http.Handler {
	return http.FileServer(http.Dir(hostdir))
}
