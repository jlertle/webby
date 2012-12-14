package webby

import (
	"net"
	"net/http"
	"strings"
)

// Set Cookie
func (w *Web) SetCookie(cookie *http.Cookie) {
	var num int

	if cookie.Path == "" {
		cookie.Path = "/"
	}

	if cookie.Domain != "" {
		goto set_cookie
	}

	cookie.Domain = w.Req.Host

	num = strings.LastIndex(cookie.Domain, "]:")
	if num != -1 {
		cookie.Domain = cookie.Domain[:num+1]
		goto skip_port_check
	}

	if cookie.Domain[len(cookie.Domain)-1] == ']' {
		goto skip_port_check
	}

	num = strings.LastIndex(cookie.Domain, ":")
	if num != -1 {
		cookie.Domain = cookie.Domain[:num]
	}

skip_port_check:

	cookie.Domain = strings.Trim(cookie.Domain, "[]")

	if net.ParseIP(cookie.Domain) != nil {
		cookie.Domain = ""
		goto set_cookie
	}

	if strings.Count(cookie.Domain, ".") <= 0 {
		cookie.Domain = ""
	}

set_cookie:
	http.SetCookie(w.Res, cookie)
}

// Get Cookie
func (w *Web) GetCookie(name string) (*http.Cookie, error) {
	return w.Req.Cookie(name)
}
