package webby

import (
	"net/http"
	"testing"
)

func TestJuntion(t *testing.T) {
	pass := FuncToRouteHandler{func(w *Web) {
		// Do nothing, because it's a pass
	}}

	fail := FuncToRouteHandler{func(w *Web) {
		t.Fail()
	}}

	w := &Web{
		Env: http.Header{},
		Req: &http.Request{
			Method: "GET",
		},
		Errors: &Errors{
			E403: fail.Function,
			E404: fail.Function,
			E500: fail.Function,
		},
	}

	jn := NewJunction().Ajax(fail).Websocket(fail).Get(pass).Post(
		fail).Delete(fail).Put(fail).Patch(fail).Options(fail).Fallback(fail).GetJunction()

	jn.View(w)

	w.Req.Method = "POST"

	jn.ALL = pass
	jn.GET = fail
	jn.POST = nil

	jn.View(w)

	jn.ALL = fail
	jn.POST = pass

	jn.View(w)

	w.Req.Method = "DELETE"

	jn.POST = fail
	jn.DELETE = pass

	jn.View(w)

	w.Req.Method = "PUT"

	jn.DELETE = fail
	jn.PUT = pass

	jn.View(w)

	w.Req.Method = "PATCH"

	jn.PUT = fail
	jn.PATCH = pass

	jn.View(w)

	w.Req.Method = "OPTIONS"

	jn.PATCH = fail
	jn.OPTIONS = pass

	jn.View(w)

	jn.OPTIONS = nil
	jn.ALL = nil

	w.Errors.E404 = pass.Function

	jn.View(w)

	w.Errors.E404 = fail.Function
	jn.OPTIONS = fail
	jn.ALL = fail

	w.Req.Method = "GET"

	w.Env.Set("X-Requested-With", "XMLHttpRequest")

	jn.AJAX = pass

	jn.View(w)

	w.Env.Set("Connection", "Upgrade")
	w.Env.Set("Upgrade", "websocket")

	jn.AJAX = fail
	jn.WS = pass

	jn.View(w)
}
