package webby

func ExampleWeb_Cmd(w *Web) {
	// Set Cmd
	w.Cmd().Set("example", func(v interface{}) interface{} {
		w.Fmt().Println(v)
		return v
	})

	// Exec Cmd
	w.Cmd().Exec("example", "Hello World")

	// Can be very useful for template engine!
}
