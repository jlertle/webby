package webby

import (
	"testing"
)

func TestRouter(t *testing.T) {
	possible_pass := FuncToRouteHandler{func(w *Web) {
		if w.Param.Get("title") != "test" {
			t.Fail()
		}

		if w.Param.GetInt("id") != 5 {
			t.Fail()
		}
	}}

	pass := FuncToRouteHandler{func(w *Web) {
		// Do nothing, it's an automactic pass!
	}}

	fail := FuncToRouteHandler{func(w *Web) {
		t.Fail()
	}}

	w := &Web{
		Param: Param{},
		pri: &webPrivate{
			path:    "/blogpost/test-5",
			curpath: "",
			cut:     false,
		},
		Errors: &Errors{
			E403: fail.Function,
			E404: fail.Function,
			E500: fail.Function,
		},
	}

	route := NewRouterHandlerMap(RouteHandlerMap{
		`^/blogpost`: NewRouterHandlerMap(RouteHandlerMap{
			`^/(?P<title>[a-z]+)-(?P<id>\d+)/?$`: possible_pass,
		}),
	})

	route.Load(w)

	w.pri.path = "/55"
	w.pri.curpath = ""

	w.Errors.E404 = pass.Function

	route.Load(w)
}
