package webby

func ExampleWeb_Xml(w *Web) {
	// Prepare structure
	data := struct {
		Title string `xml:"title,attr"`
		Life  int    `xml:"life"`
	}{}

	// Decode from Request Body
	w.Xml().DecodeReqBody(data)

	// Send it back to the client
	w.Xml().Send(data)
}
