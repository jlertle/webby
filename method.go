package webby

import (
	"reflect"
)

func execMethodInterface(w *Web, me methodInterface) {
	vc := reflect.New(reflect.Indirect(reflect.ValueOf(me)).Type())

	view := vc.MethodByName("View")
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(w)
	view.Call(in)

	in = make([]reflect.Value, 0)
	method := vc.MethodByName("Prepare")
	method.Call(in)

	if w.CutOut() {
		return
	}

	if w.Is().WebSocketRequest() {
		method = vc.MethodByName("Ws")
		method.Call(in)
		if w.CutOut() {
			goto finish
		}
	}

	if w.Is().AjaxRequest() {
		method = vc.MethodByName("Ajax")
		method.Call(in)
		if w.CutOut() {
			goto finish
		}
	}

	switch w.Req.Method {
	case "GET", "HEAD":
		method = vc.MethodByName("Get")
		method.Call(in)
	case "POST":
		method = vc.MethodByName("Post")
		method.Call(in)
	case "DELETE":
		method = vc.MethodByName("Delete")
		method.Call(in)
	case "PUT":
		method = vc.MethodByName("Put")
		method.Call(in)
	case "PATCH":
		method = vc.MethodByName("Patch")
		method.Call(in)
	case "OPTIONS":
		method = vc.MethodByName("Options")
		method.Call(in)
	}

finish:

	method = vc.MethodByName("Finish")
	method.Call(in)
}

func methodNotAllowed(w *Web) {
	w.Status = 405
	w.Fmt().Print("405 Method Now Allowed")
}

type methodInterface interface {
	View(*Web)
	Prepare()
	Ws()
	Ajax()
	Get()
	Post()
	Delete()
	Put()
	Patch()
	Options()
	Finish()
}

type Method struct {
	W *Web
}

func (me *Method) View(w *Web) {
	me.W = w
}

func (me *Method) Prepare() {
	// Do nothing
}

func (me *Method) Ws() {
	// Do nothing
}

func (me *Method) Ajax() {
	// Do nothing
}

func (me *Method) Get() {
	methodNotAllowed(me.W)
}

func (me *Method) Post() {
	methodNotAllowed(me.W)
}

func (me *Method) Delete() {
	methodNotAllowed(me.W)
}

func (me *Method) Put() {
	methodNotAllowed(me.W)
}

func (me *Method) Patch() {
	methodNotAllowed(me.W)
}

func (me *Method) Options() {
	methodNotAllowed(me.W)
}

func (me *Method) Finish() {
	// Do nothing
}
