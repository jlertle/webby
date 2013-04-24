package webby

import (
	"fmt"
	"net/http"
	"testing"
)

func TestVhost(t *testing.T) {
	pass := FuncToRouteHandler(func(w *Web) {
		// Do nothing, because it's a pass
	})

	fail := FuncToRouteHandler(func(w *Web) {
		t.Fail()
	})

	w := &Web{
		Env: http.Header{},
		Req: &http.Request{
			Host: "example.com:1234",
		},
		Errors: &Errors{
			E403: fail,
			E404: fail,
			E500: fail,
		},
		pri: &webPrivate{
			cut: false,
		},
	}

	hosts := NewVHost(VHostMap{
		`example.com`: pass,
	})

	hosts.View(w)

	w.Req.Host = "www.example.com:1234"

	w.Errors.E404 = pass

	hosts = NewVHost(VHostMap{
		`example.com`: fail,
	})

	hosts.View(w)
}

func TestVhostRegExp(t *testing.T) {
	possible_pass := FuncToRouteHandler(func(w *Web) {
		if w.Param.Get("subdomain") != "hello" {
			t.Fail()
			fmt.Println(w.Req.Host)
		}
	})

	pass := FuncToRouteHandler(func(w *Web) {
		// Do nothing, because it's a pass
	})

	fail := FuncToRouteHandler(func(w *Web) {
		t.Fail()
		fmt.Println(w.Req.Host)
	})

	w := &Web{
		Env: http.Header{},
		Req: &http.Request{
			Host: "hello.example.com:1234",
		},
		Errors: &Errors{
			E403: fail,
			E404: fail,
			E500: fail,
		},
		pri: &webPrivate{
			cut: false,
		},
		Param: Param{},
	}

	const rule = `^(?P<subdomain>[a-z]+)\.example\.com`

	vhost := NewVHostRegExp(VHostRegExpMap{
		rule: possible_pass,
	})

	vhost.View(w)

	w.Req.Host = "www.example.com:1234"
	w.Param = Param{}

	vhost = NewVHostRegExp(VHostRegExpMap{
		rule: pass,
	})

	vhost.View(w)

	w.Req.Host = "www.hello.com:1234"

	w.Errors.E404 = pass

	vhost = NewVHostRegExp(VHostRegExpMap{
		rule: fail,
	})

	vhost.View(w)
}
