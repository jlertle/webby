package webby

import (
	"net/http"
)

func ExampleWeb_Http(w *Web) {
	// Dummy Cookie
	cookie := &http.Cookie{Name: "hello"}

	// Set Cookie
	w.Http().SetCookie(cookie)

	// Dummy Handler
	handler := func(res http.ResponseWriter, req *http.Request) {
		// Do nothing
	}

	// Execute Function
	w.Http().ExecFunc(handler)
}
