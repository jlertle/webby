package webby

import (
	"reflect"
	"strings"
)

func execProtocolInterface(w *Web, pr ProtocolInterface) {
	vc := reflect.New(reflect.Indirect(reflect.ValueOf(pr)).Type())

	view := vc.MethodByName("View")
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(w)
	view.Call(in)

	in = make([]reflect.Value, 0)

	switch strings.ToLower(strings.Split(w.Req.Proto, "/")[0]) {
	case "http":
		method := vc.MethodByName("Http")
		method.Call(in)
	case "shttp", "https":
		method := vc.MethodByName("Https")
		method.Call(in)
	}
}

type ProtocolInterface interface {
	View(*Web)
	Http()
	Https()
}

type Protocol struct {
	W *Web
}

func (pr *Protocol) View(w *Web) {
	pr.W = w
}

func (pr *Protocol) Http() {
	pr.W.Error404()
}

func (pr *Protocol) Https() {
	pr.W.Error404()
}
