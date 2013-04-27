package webby

import (
	"testing"
)

func TestRouter(t *testing.T) {
	possible_pass := FuncToRouteHandler(func(w *Web) {
		if w.Param.Get("title") != "test" {
			t.Fail()
		}

		if w.Param.GetInt("id") != 5 {
			t.Fail()
		}
	})

	pass := FuncToRouteHandler(func(w *Web) {
		// Do nothing, it's an automactic pass!
	})

	fail := FuncToRouteHandler(func(w *Web) {
		t.Fail()
	})

	w := &Web{
		Param: Param{},
		pri: &webPrivate{
			path:    "/blogpost/test-5",
			curpath: "",
			cut:     false,
		},
		Errors: &Errors{
			E403: fail,
			E404: fail,
			E500: fail,
		},
	}

	route := NewRouter().RegisterHandlerMap(RouteHandlerMap{
		`^/blogpost`: NewRouter().RegisterHandlerMap(RouteHandlerMap{
			`^/(?P<title>[a-z]+)-(?P<id>\d+)/?$`: possible_pass,
		}),
	})

	route.Load(w)

	w.pri.path = "/55"
	w.pri.curpath = ""

	w.Errors.E404 = pass

	route.Load(w)
}
