package webby

import (
	"net/http"
	"net/url"
	"testing"
)

func TestVhost(t *testing.T) {
	pass := FuncToRouteHandler{func(w *Web) {
		// Do nothing, because it's a pass
	}}

	fail := FuncToRouteHandler{func(w *Web) {
		t.Fail()
	}}

	pass_router := NewBootRoute().Router(NewRouterHandlerMap(RouteHandlerMap{
		`^/$`: pass,
	})).Get()

	fail_router := NewBootRoute().Router(NewRouterHandlerMap(RouteHandlerMap{
		`^/$`: fail,
	})).Get()

	w := &Web{
		Env: http.Header{},
		Req: &http.Request{
			Host: "example.com:1234",
			URL: &url.URL{
				Path: "/",
			},
		},
		Errors: &Errors{
			E403: fail.Function,
			E404: fail.Function,
			E500: fail.Function,
		},
		pri: &webPrivate{
			path:    "/",
			curpath: "",
			cut:     false,
		},
	}

	hosts := NewVHost(VHostMap{
		`example.com`: pass_router,
	})

	hosts.View(w)

	w.Req.Host = "www.example.com:1234"

	w.Errors.E404 = pass.Function

	hosts = NewVHost(VHostMap{
		`example.com`: fail_router,
	})

	hosts.View(w)
}

func TestVhostRegExp(t *testing.T) {
	possible_pass := FuncToRouteHandler{func(w *Web) {
		if w.Param.Get("subdomain") != "hello" {
			t.Fail()
		}
	}}

	pass := FuncToRouteHandler{func(w *Web) {
		// Do nothing, because it's a pass
	}}

	fail := FuncToRouteHandler{func(w *Web) {
		t.Fail()
	}}

	possible_pass_router := NewBootRoute().Router(NewRouterHandlerMap(RouteHandlerMap{
		`^/$`: possible_pass,
	})).Get()

	pass_router := NewBootRoute().Router(NewRouterHandlerMap(RouteHandlerMap{
		`^/$`: pass,
	})).Get()

	fail_router := NewBootRoute().Router(NewRouterHandlerMap(RouteHandlerMap{
		`^/$`: fail,
	})).Get()

	w := &Web{
		Env: http.Header{},
		Req: &http.Request{
			Host: "hello.example.com:1234",
			URL: &url.URL{
				Path: "/",
			},
		},
		Errors: &Errors{
			E403: fail.Function,
			E404: fail.Function,
			E500: fail.Function,
		},
		pri: &webPrivate{
			path:    "/",
			curpath: "",
			cut:     false,
		},
		Param: Param{},
	}

	const rule = `^(?P<subdomain>[a-z]+)\.example\.com`

	vhost := NewVHostRegExp(VHostRegExpMap{
		rule: possible_pass_router,
	})

	vhost.View(w)

	w.Req.Host = "www.example.com:1234"
	w.Param = Param{}

	vhost = NewVHostRegExp(VHostRegExpMap{
		rule: pass_router,
	})

	vhost.View(w)

	w.Req.Host = "www.hello.com:1234"

	w.Errors.E404 = pass.Function

	vhost = NewVHostRegExp(VHostRegExpMap{
		rule: fail_router,
	})

	vhost.View(w)
}
