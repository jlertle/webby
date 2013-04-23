package webby

func ExampleWeb_Json(w *Web) {
	// Prepare structure
	data := struct {
		Title string `json:"title"`
		Life  int    `json:"life"`
	}{}

	// Decode from Request Body
	w.Json().DecodeReqBody(data)

	// Send it back to the client
	w.Json().Send(data)
}
