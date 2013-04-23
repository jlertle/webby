package webby

func ExampleWeb_Html(w *Web) {
	// Prepare Struct
	data := struct {
		A, B string
	}{
		A: "Hello",
		B: "World",
	}

	// Prepare Html Template
	const htmlstr = "<h1>{{.A}} {{.B}}</h1>"

	// Render and Send to Client, don't worry it won't send the headers too early, like it does in wordpress! (That another story) I have made sure of that!
	w.Html().RenderSend(htmlstr, data)

	// Same as above, except it's grabs the template out of the file.
	w.Html().RenderFileSend("./example.html", data)

	// Something more advanced!
	// Let say the file content is `{{define "index"}}<h1>{{.A}} {{.B}}</h1>{{end}}`

	// The Bootstrap is the best place to set that!
	w.Html().SetDefaultFiles("./file.html")

	// Render "index" and send!
	w.Html().DefaultRenderSend("index", data)

	// If not in debug mode the template file will automatically cache to ram!
}
