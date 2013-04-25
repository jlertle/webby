package webby

import (
	"testing"
)

type MiddlewareDummy struct {
	Middleware
}

func (mid *MiddlewareDummy) Pre() {
	mid.W.Param.Set("result", "PRE")
}

func (mid *MiddlewareDummy) Post() {
	mid.W.Param.Set("result", "POST")
}

func TestMiddleware(t *testing.T) {
	w := &Web{
		Param: Param{},
		pri: &webPrivate{
			cut: false,
		},
	}

	result := func() string {
		return w.Param.Get("result")
	}

	mid := NewMiddlewares().Register(&MiddlewareDummy{}).Init(w)
	mid.Pre()

	if result() != "PRE" {
		t.Fail()
	}

	mid.Post()

	if result() != "POST" {
		t.Fail()
	}
}
