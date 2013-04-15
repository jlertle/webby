package webby

import (
	"encoding/json"
)

type Json struct {
	w *Web
}

func (w *Web) Json() Json {
	return Json{w}
}

// Output in JSON
func (j Json) Send(v interface{}) {
	json.NewEncoder(j.w).Encode(v)
}
