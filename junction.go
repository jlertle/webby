package webby

type Junction struct {
	ALL, GET, POST, DELETE, PUT, AJAX, WS RouteHandler
}

func (jn Junction) View(w *Web) {
	if w.IsWebSocketRequest() {
		if jn.WS != nil {
			jn.WS.View(w)
			return
		}
	}

	if w.IsAjaxRequest() {
		if jn.AJAX != nil {
			jn.AJAX.View(w)
			return
		}
	}

	switch w.Req.Method {
	case "GET":
		if jn.GET != nil {
			jn.GET.View(w)
			return
		}
	case "POST":
		if jn.POST != nil {
			jn.POST.View(w)
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
	}

	if jn.ALL != nil {
		jn.ALL.View(w)
		return
	}

	w.Error404()
	return
}
