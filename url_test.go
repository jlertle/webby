package webby

func ExampleWeb_Url(w *Web) {
	// Get Absolute Url
	url := w.Url().Absolute("/")

	// Get Absolute Url (Https)
	url = w.Url().AbsoluteHttps("/")

	// Output Url to Client
	w.Fmt().Println(url)
}
