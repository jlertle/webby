package webby

import (
	"net/http"
)

// Create new File Server and returns RouteHandler
func FileServer(dir string) RouteHandler {
	adir := dir
	return FuncToRouteHandler(func(w *Web) {
		w.Http().Exec(http.FileServer(http.Dir(adir)))
	})
}
