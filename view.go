package webby

// Overridable Error403 Function
//
// Note:  Overriding is useful for custom 403 page
var Error403 = func(w *Web) {
	w.Status = 403
	w.Print("<h1>403 Forbidded</h1>")
}

func (web *Web) Error403() {
	Error403(web)
}

// Overridable Error404 Function
//
// Note:  Overriding is useful for custom 404 page
var Error404 = func(w *Web) {
	w.Status = 404
	w.Print("<h1>404 Not Found</h1>")
}

func (web *Web) Error404() {
	Error404(web)
}

// Overridable Error500 Function
//
// Note:  Overriding is useful for custom 500 page
var Error500 = func(w *Web) {
	w.Status = 500
	w.Print("<h1>500 Internal Server Error</h1>")
}

func (web *Web) Error500() {
	Error500(web)
}

// Default Index View
func index(w *Web) {
	w.Print("<h1>Hello World!</h1>")
}

// Push Index View to Router
func init() {
	Route.Register("^/$", index)
}
