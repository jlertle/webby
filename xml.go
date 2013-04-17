package webby

import (
	"encoding/xml"
	"io"
)

type Xml struct {
	w *Web
}

func (w *Web) Xml() Xml {
	return Xml{w}
}

// Shortcut to encoding/xml.Marshal
func (x Xml) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Shortcut to encoding/xml.MarshalIndent
func (x Xml) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return xml.MarshalIndent(v, prefix, indent)
}

// Shortcut to encoding/xml.Unmarshal
func (x Xml) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// Shortcut to encoding/xml.NewDecoder
func (x Xml) NewDecoder(r io.Reader) *xml.Decoder {
	return xml.NewDecoder(r)
}

// Shortcut to encoding/xml.NewEncoder
func (x Xml) NewEncoder(w io.Writer) *xml.Encoder {
	return xml.NewEncoder(w)
}

// Output in XML
func (x Xml) Send(v interface{}) {
	xml.NewEncoder(x.w).Encode(v)
}

// Decode Request Body
func (x Xml) DecodeReqBody(v interface{}) error {
	return x.NewDecoder(x.w.Req.Body).Decode(v)
}
