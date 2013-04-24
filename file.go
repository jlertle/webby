package webby

import (
	"net/http"
)

// Create new File Server and returns RouteHandler
func FileServer(dir string) RouteHandler {
	adir := dir
	return FuncToRouteHandler(func(w *Web) {
		http.StripPrefix(w.pri.curpath, http.FileServer(http.Dir(adir))).ServeHTTP(w, w.Req)
	})
}
