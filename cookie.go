package webby

import (
	"net"
	"net/http"
	"strings"
	"time"
)

func init() {
	HtmlFuncBoot.Register(func(w *Web) {
		// Get Cookie Value
		w.HtmlFunc["cookie"] = func(name string) string {
			cookie, err := w.Cookie(name).Get()
			if err != nil {
				return ""
			}
			return cookie.Value
		}
	})
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

// Cookie
func (w *Web) Cookie(name string) PipeCookie {
	return NewCookie(w, name)
}

// Set Value
func (c PipeCookie) Value(value string) PipeCookie {
	c.c.Value = value
	return c
}

// Set Path
func (c PipeCookie) Path(path string) PipeCookie {
	c.c.Path = path
	return c
}

// Set Domain
func (c PipeCookie) Domain(domain string) PipeCookie {
	c.c.Domain = domain
	return c
}

// Set Expiry Time of Cookie.
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

// Make Cookie Secure
func (c PipeCookie) Secure() PipeCookie {
	c.c.Secure = true
	return c
}

// Make Cookie Http Only
func (c PipeCookie) HttpOnly() PipeCookie {
	c.c.HttpOnly = true
	return c
}

// Get *http.Cookie, if Value is not set it will try to get the Cookie from the User Request!
func (c PipeCookie) Get() (*http.Cookie, error) {
	if c.c.Value != "" {
		return c.c, nil
	}
	return c.w.Req.Cookie(c.c.Name)
}

// Delete Cookie
func (c PipeCookie) Delete() PipeCookie {
	return c.Value("Delete-Me").MaxAge(-1).SaveRes()
}

// Save (Set) Cookie to Response
func (c PipeCookie) SaveRes() PipeCookie {
	http.SetCookie(c.w, c.pre(c.c))
	return c
}

// Prepare Cookie
func (c PipeCookie) pre(cookie *http.Cookie) *http.Cookie {
	var num int
	w := c.w

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

// Save (Add) Cookie to Request, It won't send anything out to the client.
// But it is a useful feature for CSRF protection for example!.
func (c PipeCookie) SaveReq() PipeCookie {
	c.w.Req.AddCookie(c.c)
	return c
}

// Set to Expire after an hour
func (c PipeCookie) Hour() PipeCookie {
	return c.Expires(time.Now().Add(1 * time.Hour))
}

// Set to Expire after 6 Hours
func (c PipeCookie) Hour6() PipeCookie {
	return c.Expires(time.Now().Add(6 * time.Hour))
}

// Set to Expire after 12 Hours
func (c PipeCookie) Hour12() PipeCookie {
	return c.Expires(time.Now().Add(12 * time.Hour))
}

// Set to Expire after 1 Day
func (c PipeCookie) Day() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 0, 1))
}

// Set to Expire after 1 Week
func (c PipeCookie) Week() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 0, 1*7))
}

// Set to Expire after 2 Week
func (c PipeCookie) Week2() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 0, 2*7))
}

// Set to Expire after 1 Month
func (c PipeCookie) Month() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 1, 0))
}

// Set to Expire after 3 Month
func (c PipeCookie) Month3() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 3, 0))
}

// Set to Expire after 6 Month
func (c PipeCookie) Month6() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 6, 0))
}

// Set to Expire after 9 Month
func (c PipeCookie) Month9() PipeCookie {
	return c.Expires(time.Now().AddDate(0, 9, 0))
}

// Set to Expire after 1 Year
func (c PipeCookie) Year() PipeCookie {
	return c.Expires(time.Now().AddDate(1, 0, 0))
}
