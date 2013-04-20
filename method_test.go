package webby

import (
	"net/http"
	"testing"
)

type MethodDummy struct {
	Method
}

func (me *MethodDummy) Prepare() {
	me.W.Param.Set("prepare", "PREPARE")
}

func (me *MethodDummy) Ws() {
	me.W.Param.Set("ws", "WS")
}

func (me *MethodDummy) Ajax() {
	me.W.Param.Set("ajax", "AJAX")
}

func (me *MethodDummy) Finish() {
	me.W.Param.Set("finish", "FINISH")
}

func (me *MethodDummy) Get() {
	me.W.Param.Set("method", "GET")
}

func (me *MethodDummy) Post() {
	me.W.Param.Set("method", "POST")
}

func (me *MethodDummy) Delete() {
	me.W.Param.Set("method", "DELETE")
}

func (me *MethodDummy) Put() {
	me.W.Param.Set("method", "PUT")
}

func (me *MethodDummy) Patch() {
	me.W.Param.Set("method", "PATCH")
}

func (me *MethodDummy) Options() {
	me.W.Param.Set("method", "OPTIONS")
}

func TestMethod(t *testing.T) {
	w := &Web{
		Param: Param{},
		Env:   http.Header{},
		Req: &http.Request{
			Method: "GET",
		},
		pri: &webPrivate{
			cut: false,
		},
	}

	prepare := func() string {
		return w.Param.Get("prepare")
	}

	ws := func() string {
		return w.Param.Get("ws")
	}

	ajax := func() string {
		return w.Param.Get("ajax")
	}

	finish := func() string {
		return w.Param.Get("finish")
	}

	method := func() string {
		return w.Param.Get("method")
	}

	w.RouteDealer(&MethodDummy{})

	if prepare() != "PREPARE" {
		t.Fail()
	}

	if ws() == "WS" {
		t.Fail()
	}

	if ajax() == "AJAX" {
		t.Fail()
	}

	if finish() != "FINISH" {
		t.Fail()
	}

	if method() != "GET" {
		t.Fail()
	}

	slices := []string{"POST", "DELETE", "PUT", "PATCH", "OPTIONS"}

	for _, value := range slices {
		w.Req.Method = value
		w.RouteDealer(&MethodDummy{})
		if method() != value {
			t.Fail()
		}
	}
}
