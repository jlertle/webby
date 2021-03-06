package webby

import (
	"io/ioutil"
	"net"
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

func (w *Web) initSecure() {
	if w.Env.Get("X-Secure-Mode") != "" {
		w.Req.Proto = "S" + w.Req.Proto
		w.Env.Del("X-Secure-Mode")
	}
}

// Get Remote Address (IP Address) without port number!
func (w *Web) RemoteAddr() string {
	ip, _, _ := net.SplitHostPort(w.Req.RemoteAddr)
	return ip
}
