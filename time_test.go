package webby

func ExampleWeb_Time(w *Web) {
	// Get current time
	curtime := w.Time().Now()

	// Output to client
	w.Fmt().Println(curtime)

	// Set Timezone on user request level
	w.Time().SetZone("Europe/London")
}
