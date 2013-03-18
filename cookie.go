package webby

import (
	"net"
	"net/http"
	"strings"
	"time"
)

// Prepare Cookie
func (w *Web) preCookie(cookie *http.Cookie) *http.Cookie {
	var num int

	if cookie.Path == "" {
		cookie.Path = "/"
	}

	if cookie.Domain != "" {
		goto release_cookie
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
		goto release_cookie
	}

	if strings.Count(cookie.Domain, ".") <= 0 {
		cookie.Domain = ""
	}

release_cookie:

	return cookie
}

// Set Cookie
func (w *Web) SetCookie(cookie *http.Cookie) {
	http.SetCookie(w, w.preCookie(cookie))
}

// Get Cookie
func (w *Web) GetCookie(name string) (*http.Cookie, error) {
	return w.Req.Cookie(name)
}

// Delete Cookie
func (w *Web) DeleteCookie(name string) {
	w.Cookie(name).Value("Delete-Me").MaxAge(-1).SaveRes()
}

// Chainable version of 'net/http.Cookie'
type PipeCookie struct {
	w *Web
	c *http.Cookie
}

// New Cookie
func NewCookie(w *Web, name string) PipeCookie {
	return PipeCookie{
		w: w,
		c: &http.Cookie{Name: name},
	}
}

// Alias of New Cookie
func (w *Web) Cookie(name string) PipeCookie {
	return NewCookie(w, name)
}

// New Cookie
func (w *Web) NewCookie(name string) PipeCookie {
	return NewCookie(w, name)
}

func (c PipeCookie) Value(value string) PipeCookie {
	c.c.Value = value
	return c
}

func (c PipeCookie) Path(path string) PipeCookie {
	c.c.Path = path
	return c
}

func (c PipeCookie) Domain(domain string) PipeCookie {
	c.c.Domain = domain
	return c
}

func (c PipeCookie) Expires(expires time.Time) PipeCookie {
	c.c.Expires = expires
	return c
}

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
func (c PipeCookie) MaxAge(maxage int) PipeCookie {
	c.c.MaxAge = maxage
	return c
}

func (c PipeCookie) Secure() PipeCookie {
	c.c.Secure = true
	return c
}

func (c PipeCookie) HttpOnly() PipeCookie {
	c.c.HttpOnly = true
	return c
}

func (c PipeCookie) Get() *http.Cookie {
	return c.c
}

// Save (Set) Cookie to Response
func (c PipeCookie) SaveRes() PipeCookie {
	c.w.SetCookie(c.c)
	return c
}

// Save (Add) Cookie to Request, It won't send anything out to the client.
// But it is a useful feature for CSRF protection for example!.
func (c PipeCookie) SaveReq() PipeCookie {
	c.w.Req.AddCookie(c.c)
	return c
}
