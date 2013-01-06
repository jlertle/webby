package webby

import (
	"net/http"
)

type HttpRouteHandler struct {
	http.Handler
}

func (ht HttpRouteHandler) View(w *Web) {
	ht.ServeHTTP(w, w.Req)
}
