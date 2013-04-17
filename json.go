package webby

import (
	"encoding/json"
	"io"
)

type Json struct {
	w *Web
}

func (w *Web) Json() Json {
	return Json{w}
}

// Shortcut to encoding/json.Marshal
func (j Json) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Shortcut to encoding/json.MarshalIndent
func (j Json) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Shortcut to encoding/json.Unmarshal
func (j Json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Shortcut to encoding/json.NewDecoder
func (j Json) NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

// Shortcut to encoding/json.NewEncoder
func (j Json) NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}

// Send Json output to client.
func (j Json) Send(v interface{}) {
	j.NewEncoder(j.w).Encode(v)
}

// Decode Request Body
func (j Json) DecodeReqBody(v interface{}) error {
	return j.NewDecoder(j.w.Req.Body).Decode(v)
}
