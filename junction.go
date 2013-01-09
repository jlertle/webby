package webby

import (
	"strings"
)

type Junction struct {
	ALL, GET, POST, HEAD, DELETE, PUT, PATCH, OPTIONS, AJAX, WS RouteHandler
}

func (jn Junction) View(w *Web) {
	if w.IsWebSocketRequest() {
		if jn.WS != nil {
			jn.WS.View(w)
			return
		}
	}

	switch w.Req.Method {
	case "GET":
		if w.IsAjaxRequest() {
			if jn.AJAX != nil {
				jn.AJAX.View(w)
				return
			}
		}
		if jn.GET != nil {
			jn.GET.View(w)
			return
		}
	case "POST":
		if w.IsAjaxRequest() {
			if jn.AJAX != nil {
				jn.AJAX.View(w)
				return
			}
		}
		if jn.POST != nil {
			jn.POST.View(w)
			return
		}
	case "HEAD":
		if jn.HEAD != nil {
			jn.HEAD.View(w)
			return
		}
	case "DELETE":
		if jn.DELETE != nil {
			jn.DELETE.View(w)
			return
		}
	case "PUT":
		if jn.PUT != nil {
			jn.PUT.View(w)
			return
		}
	case "PATCH":
		if jn.PATCH != nil {
			jn.PATCH.View(w)
			return
		}
	case "OPTIONS":
		if jn.OPTIONS != nil {
			jn.OPTIONS.View(w)
			return
		}
	}

	if jn.ALL != nil {
		jn.ALL.View(w)
		return
	}

	w.Error404()
	return
}

type Protocol struct {
	ALL, HTTP, HTTPS RouteHandler
}

func (pr Protocol) View(w *Web) {
	switch strings.ToLower(strings.Split(w.Req.Proto, "/")[0]) {
	case "http":
		if pr.HTTP != nil {
			pr.HTTP.View(w)
			return
		}
	case "shttp":
		if pr.HTTPS != nil {
			pr.HTTPS.View(w)
			return
		}
	}

	if pr.ALL != nil {
		pr.ALL.View(w)
		return
	}

	w.Error404()
	return
}
