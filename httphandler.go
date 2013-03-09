package webby

import (
	"net/http"
)

// 'net/http.Handler' Adapter.  Implement RouterHandler interface 
type HttpRouteHandler struct {
	http.Handler
}

func (ht HttpRouteHandler) View(w *Web) {
	ht.ServeHTTP(w, w.Req)
}
