package webby

// Demuxer for Request Method. Implement RouteHandler interface!
type Junction struct {
	ALL, GET, POST, DELETE, PUT, PATCH, OPTIONS, AJAX, WS RouteHandler
}

func (jn Junction) View(w *Web) {
	if w.IsWebSocketRequest() {
		if jn.WS != nil {
			jn.WS.View(w)
			return
		}
	}

	switch w.Req.Method {
	case "GET", "HEAD":
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

// Chainable version of Junction
type PipeJunction struct {
	jn Junction
}

// New Junction constructor
func NewJunction() PipeJunction {
	return PipeJunction{Junction{}}
}

// Get Junction
func (pi PipeJunction) GetJunction() Junction {
	return pi.jn
}

// Set Get
func (pi PipeJunction) Get(get RouteHandler) PipeJunction {
	pi.jn.GET = get
	return pi
}

// Set Post
func (pi PipeJunction) Post(post RouteHandler) PipeJunction {
	pi.jn.POST = post
	return pi
}

// Set Delete
func (pi PipeJunction) Delete(del RouteHandler) PipeJunction {
	pi.jn.DELETE = del
	return pi
}

// Set Put
func (pi PipeJunction) Put(put RouteHandler) PipeJunction {
	pi.jn.PUT = put
	return pi
}

// Set Patch
func (pi PipeJunction) Patch(patch RouteHandler) PipeJunction {
	pi.jn.PATCH = patch
	return pi
}

// Set Options
func (pi PipeJunction) Options(options RouteHandler) PipeJunction {
	pi.jn.OPTIONS = options
	return pi
}

// Set Ajax
func (pi PipeJunction) Ajax(ajax RouteHandler) PipeJunction {
	pi.jn.AJAX = ajax
	return pi
}

// Set Websocket
func (pi PipeJunction) Websocket(ws RouteHandler) PipeJunction {
	pi.jn.WS = ws
	return pi
}

// Set All
func (pi PipeJunction) All(all RouteHandler) PipeJunction {
	pi.jn.ALL = all
	return pi
}

// Alais of All
func (pi PipeJunction) Any(any RouteHandler) PipeJunction {
	return pi.All(any)
}

// Alais of All
func (pi PipeJunction) Fallback(fallback RouteHandler) PipeJunction {
	return pi.All(fallback)
}
