package webby

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Convert String to int64
func toInt(number string) (int64, error) {
	return strconv.ParseInt(number, 10, 64)
}

// Convert String to Uint64
func toUint(number string) (uint64, error) {
	return strconv.ParseUint(number, 10, 64)
}

// Convert String to float64
func toFloat(number string) (float64, error) {
	return strconv.ParseFloat(number, 64)
}

func (w *Web) initWriter() {
	if w.Req.Method == "HEAD" {
		w.pri.reswrite = ioutil.Discard
		w.Header().Set("Connection", "close")
		return
	}
	w.pri.reswrite = w.web
	w.Header().Set("Content-Encoding", "plain")
}

func (w *Web) initTrueHost() {
	switch {
	case w.Env.Get("Host") != "":
		w.Req.Host = w.Env.Get("Host")
	case w.Env.Get("X-Forwarded-Host") != "":
		w.Req.Host = w.Env.Get("X-Forwarded-Host")
	case w.Env.Get("X-Forwarded-Server") != "":
		w.Req.Host = w.Env.Get("X-Forwarded-Server")
	}
	w.Req.URL.Host = w.Req.Host
}

func (w *Web) initTrueRemoteAddr() {
	switch {
	case w.Env.Get("X-Real-Ip") != "":
		w.Req.RemoteAddr = w.Env.Get("X-Real-Ip") + ":1234"
		return
	case w.Env.Get("X-Forwarded-For") != "":
		w.Req.RemoteAddr = w.Env.Get("X-Forwarded-For") + ":1234"
		return
	}
}

func (w *Web) initTruePath() {
	switch {
	case w.Env.Get("X-Original-Url") != "":
		// For compatibility with IIS
		urls := strings.Split(w.Env.Get("X-Original-Url"), "?")
		w.Req.URL.Path = urls[0]
		w.pri.path = w.Req.URL.Path

		if len(urls) < 2 {
			return
		}

		w.Req.URL.RawQuery = strings.Join(urls[1:], "?")
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

var stripPortFromAddr = regexp.MustCompile("^(.*)(:(\\d+))$")

// Get Remote Address (IP Address) without port number!
func (w *Web) RemoteAddr() string {
	return stripPortFromAddr.FindStringSubmatch(w.Req.RemoteAddr)[1]
}

// Output in JSON
func (w *Web) JsonSend(v interface{}) {
	json.NewEncoder(w).Encode(v)
}

// Output in XML
func (w *Web) XmlSend(v interface{}) {
	xml.NewEncoder(w).Encode(v)
}
