package webby

import (
	"strings"
)

type Is struct {
	w *Web
}

func (w *Web) Is() Is {
	return Is{w}
}

// Is Ajax Request
func (i Is) AjaxRequest() bool {
	return i.w.Env.Get("X-Requested-With") == "XMLHttpRequest"
}

// Is WebSocket Request
func (i Is) WebSocketRequest() bool {
	return i.w.Env.Get("Connection") == "Upgrade" && i.w.Env.Get("Upgrade") == "websocket"
}

// Is Do Not Track
func (i Is) DNT() bool {
	return i.w.Env.Get("Dnt") == "1" || i.w.Env.Get("X-Do-Not-Track") == "1"
}

// Is Secure
func (i Is) Secure() bool {
	return strings.ToLower(strings.Split(i.w.Req.Proto, "/")[0]) == "shttp"
}
