package webby

// Default Index View
func index(w *Web) {
	w.Fmt().Print("<h1>Hello World!</h1>")
}

// Push Index View to Router
func init() {
	Route.Register("^/$", index)
}
