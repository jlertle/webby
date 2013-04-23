package webby

func ExampleWeb_Fmt(w *Web) {
	// Print to Client
	w.Fmt().Print("Hello", " World")

	// Print to Client with New Line
	w.Fmt().Println("Hello World")

	// Print to Client (Format)
	w.Fmt().Printf("%s %s", "Hello", "World")
}
