package webby

// Overridable Error403 Function
//
// Note:  Overriding is useful for custom 403 page
var Error403 = func(w *Web) {
	w.Print("<h1>403 Forbidden</h1>")
}

func (w *Web) Error403() {
	w.Status = 403
	Error403(w)
}

// Overridable Error404 Function
//
// Note:  Overriding is useful for custom 404 page
var Error404 = func(w *Web) {
	w.Print("<h1>404 Not Found</h1>")
}

func (w *Web) Error404() {
	w.Status = 404
	Error404(w)
}

// Overridable Error500 Function
//
// Note:  Overriding is useful for custom 500 page
var Error500 = func(w *Web) {
	w.Print("<h1>500 Internal Server Error</h1>")
}

func (w *Web) Error500() {
	w.Status = 500
	Error500(w)
}

// Default Index View
func index(w *Web) {
	w.Print("<h1>Hello World!</h1>")
}

// Push Index View to Router
func init() {
	Route.Register("^/$", index)
}
