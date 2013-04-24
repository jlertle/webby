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
type Cookie struct {
	w *Web
	c *http.Cookie
}

// New Cookie
func NewCookie(w *Web, name string) Cookie {
	return Cookie{
		w: w,
		c: &http.Cookie{Name: name},
	}
}

// Cookie
func (w *Web) Cookie(name string) Cookie {
	return NewCookie(w, name)
}

// Set Value
func (c Cookie) Value(value string) Cookie {
	c.c.Value = value
	return c
}

// Set Path
func (c Cookie) Path(path string) Cookie {
	c.c.Path = path
	return c
}

// Set Domain
func (c Cookie) Domain(domain string) Cookie {
	c.c.Domain = domain
	return c
}

// Set Expiry Time of Cookie.
func (c Cookie) Expires(expires time.Time) Cookie {
	c.c.Expires = expires
	return c
}

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
func (c Cookie) MaxAge(maxage int) Cookie {
	c.c.MaxAge = maxage
	return c
}

// Make Cookie Secure
func (c Cookie) Secure() Cookie {
	c.c.Secure = true
	return c
}

// Make Cookie Http Only
func (c Cookie) HttpOnly() Cookie {
	c.c.HttpOnly = true
	return c
}

// Get *http.Cookie, if Value is not set it will try to get the Cookie from the User Request!
func (c Cookie) Get() (*http.Cookie, error) {
	if c.c.Value != "" {
		return c.c, nil
	}
	return c.w.Req.Cookie(c.c.Name)
}

// Delete Cookie
func (c Cookie) Delete() Cookie {
	return c.Value("Delete-Me").MaxAge(-1).SaveRes()
}

// Save (Set) Cookie to Response
func (c Cookie) SaveRes() Cookie {
	http.SetCookie(c.w, c.pre(c.c))
	return c
}

// Prepare Cookie
func (c Cookie) pre(cookie *http.Cookie) *http.Cookie {
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
func (c Cookie) SaveReq() Cookie {
	c.w.Req.AddCookie(c.c)
	return c
}

// Set to Expire after an hour
func (c Cookie) Hour() Cookie {
	return c.Expires(time.Now().Add(1 * time.Hour))
}

// Set to Expire after 6 Hours
func (c Cookie) Hour6() Cookie {
	return c.Expires(time.Now().Add(6 * time.Hour))
}

// Set to Expire after 12 Hours
func (c Cookie) Hour12() Cookie {
	return c.Expires(time.Now().Add(12 * time.Hour))
}

// Set to Expire after 1 Day
func (c Cookie) Day() Cookie {
	return c.Expires(time.Now().AddDate(0, 0, 1))
}

// Set to Expire after 1 Week
func (c Cookie) Week() Cookie {
	return c.Expires(time.Now().AddDate(0, 0, 1*7))
}

// Set to Expire after 2 Week
func (c Cookie) Week2() Cookie {
	return c.Expires(time.Now().AddDate(0, 0, 2*7))
}

// Set to Expire after 1 Month
func (c Cookie) Month() Cookie {
	return c.Expires(time.Now().AddDate(0, 1, 0))
}

// Set to Expire after 3 Month
func (c Cookie) Month3() Cookie {
	return c.Expires(time.Now().AddDate(0, 3, 0))
}

// Set to Expire after 6 Month
func (c Cookie) Month6() Cookie {
	return c.Expires(time.Now().AddDate(0, 6, 0))
}

// Set to Expire after 9 Month
func (c Cookie) Month9() Cookie {
	return c.Expires(time.Now().AddDate(0, 9, 0))
}

// Set to Expire after 1 Year
func (c Cookie) Year() Cookie {
	return c.Expires(time.Now().AddDate(1, 0, 0))
}
