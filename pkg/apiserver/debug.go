//+build !release

package apiserver

import (
	"expvar"
	"log"
)

func init() {
	Initializer(3, func(svr *ApiServer) error {
		log.Println("initialize expvar")

		svr.afterRouter(func(*ApiServer) {
			svr.Handle("/debug/vars", expvar.Handler())
		})
		return nil
	})
}
