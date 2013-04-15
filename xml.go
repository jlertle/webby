package webby

import (
	"encoding/xml"
)

type Xml struct {
	w *Web
}

func (w *Web) Xml() Xml {
	return Xml{w}
}

// Output in XML
func (x Xml) Send(v interface{}) {
	xml.NewEncoder(x.w).Encode(v)
}
