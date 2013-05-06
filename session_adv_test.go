package webby

func ExampleWeb_SessionAdv(w *Web) {
	// Set Session
	w.SessionAdv().Set("key", "value")

	// Save Session
	w.SessionAdv().Save()

	// Get Session
	str := w.SessionAdv().Get("key").(string)

	// This should output 'value' to client
	w.Fmt().Print(str)
}
