package webby

import (
	"net/http"
)

func FileServer(dir string) RouteHandler {
	adir := dir
	return FuncToRouteHandler{func(w *Web) {
		http.StripPrefix(w.pri.curpath, http.FileServer(http.Dir(adir))).ServeHTTP(w, w.Req)
	}}
}
