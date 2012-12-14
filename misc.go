package webby

import (
	"regexp"
	"strconv"
)

// Convert String to int64
func toInt(number string) (int64, error) {
	return strconv.ParseInt(number, 10, 64)
}

func (w *Web) initTrueHost() {
	switch {
	case w.Env.Get("X-Forwarded-Host") != "":
		w.Req.Host = w.Env.Get("X-Forwarded-Host")
		w.Req.URL.Host = w.Env.Get("X-Forwarded-Host")
		return
	case w.Env.Get("X-Forwarded-Server") != "":
		w.Req.Host = w.Env.Get("X-Forwarded-Server")
		w.Req.URL.Host = w.Env.Get("X-Forwarded-Server")
		return
	}

	w.Req.URL.Host = w.Req.Host
}

func (w *Web) initTrueRemoteAddr() {
	if w.Env.Get("X-Forwarded-For") != "" {
		w.Req.RemoteAddr = w.Env.Get("X-Forwarded-For") + ":1234"
		return
	}
}

// Get Absolute URL, you can leave relative_url blank just to get the root url.
func (w *Web) AbsoluteUrl(relative_url string) string {
	if w.Req.URL.Host != "" {
		return "http://" + w.Req.URL.Host + relative_url
	}

	return relative_url
}

// Get Absolute URL (https), you can leave relative_url blank just to get the root url.
func (w *Web) AbsoluteUrlHttps(relative_url string) string {
	if w.Req.URL.Host != "" {
		return "https://" + w.Req.URL.Host + relative_url
	}

	return relative_url
}

// Is Ajax Request
func (w *Web) IsAjaxRequest() bool {
	return w.Env.Get("X-Requested-With") == "XMLHttpRequest"
}

// Is WebSocket Request
func (w *Web) IsWebSocketRequest() bool {
	return w.Env.Get("Connection") == "Upgrade" && w.Env.Get("Upgrade") == "websocket"
}

// Is Do Not Track
func (w *Web) IsDNT() bool {
	return w.Env.Get("Dnt") == "1" || w.Env.Get("X-Do-Not-Track") == "1"
}

// Html Form Memory Limit
var FormMemoryLimit = int64(16 * 1024 * 1024)

// Parse Form
func (w *Web) ParseForm() error {
	return w.Req.ParseMultipartForm(FormMemoryLimit)
}

var stripPortFromAddr = regexp.MustCompile("^(.*)(:(\\d+))$")

// Get Remote Address (IP Address) without port number!
func (w *Web) RemoteAddr() string {
	return stripPortFromAddr.FindStringSubmatch(w.Req.RemoteAddr)[1]
}
