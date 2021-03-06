package webby

import (
	"fmt"
	"sync"
)

// URL Reverse Map
type URLReverseMap map[string]string

// URL Reverse Data Type
type URLReverse struct {
	sync.RWMutex
	urls URLReverseMap
}

func (u *URLReverse) register(name, format string) {
	u.Lock()
	defer u.Unlock()
	u.urls[name] = format
}

func (u *URLReverse) initUrls() {
	u.Lock()
	defer u.Unlock()
	if u.urls == nil {
		u.urls = URLReverseMap{}
	}
}

// Url Reverse Register
func (u *URLReverse) Register(name, format string) *URLReverse {
	u.initUrls()
	u.register(name, format)
	return u
}

// Url Reverse Register Map
func (u *URLReverse) RegisterMap(urls URLReverseMap) *URLReverse {
	if urls == nil {
		return u
	}

	u.initUrls()

	for name, format := range urls {
		u.register(name, format)
	}

	return u
}

// Print relative URL to string
func (u *URLReverse) Print(name string, a ...interface{}) string {
	u.RLock()
	defer u.RUnlock()
	return fmt.Sprintf(u.urls[name], a...)
}

// Default URL Reverse!
var URLRev = &URLReverse{}

type Url struct {
	w *Web
}

func (w *Web) Url() Url {
	return Url{w}
}

// Get Absolute URL, you can leave relative_url blank just to get the root url.
func (u Url) Absolute(relative_url string) string {
	if u.w.Is().Secure() {
		return u.AbsoluteHttps(relative_url)
	}
	return u.AbsoluteHttp(relative_url)
}

// Get Absolute URL (http), you can leave relative_url blank just to get the root url.
func (u Url) AbsoluteHttp(relative_url string) string {
	w := u.w
	if w.Req.URL.Host != "" {
		return "http://" + w.Req.URL.Host + relative_url
	}

	return relative_url
}

// Get Absolute URL (https), you can leave relative_url blank just to get the root url.
func (u Url) AbsoluteHttps(relative_url string) string {
	w := u.w
	if w.Req.URL.Host != "" {
		return "https://" + w.Req.URL.Host + relative_url
	}

	return relative_url
}

// Shortcut to Reverse
func (u Url) Reverse(name string, a ...interface{}) string {
	return URLRev.Print(name, a...)
}

func (u Url) code301() int {
	if u.w.Status == 200 {
		return 301
	}
	return u.w.Status
}

func (u Url) code303() int {
	if u.w.Status == 200 {
		return 303
	}
	return u.w.Status
}

// Convert current path to Https
func (u Url) ToHttps() {
	defer u.w.WriteHeader(u.code301())
	u.w.Header().Set("Location", u.AbsoluteHttps(u.w.Req.URL.Path))
}

// Convert current path to Http
func (u Url) ToHttp() {
	defer u.w.WriteHeader(u.code301())
	u.w.Header().Set("Location", u.AbsoluteHttp(u.w.Req.URL.Path))
}

// Redirect client to relative_url
func (u Url) Redirect(relative_url string) {
	defer u.w.WriteHeader(u.code303())
	u.w.Header().Set("Location", u.Absolute(relative_url))
}

// Redirect client to relative_url (Http Only)
func (u Url) RedirectHttp(relative_url string) {
	defer u.w.WriteHeader(u.code303())
	u.w.Header().Set("Location", u.AbsoluteHttp(relative_url))
}

// Redirect client to relative_url (Https Only)
func (u Url) RedirectHttps(relative_url string) {
	defer u.w.WriteHeader(u.code303())
	u.w.Header().Set("Location", u.AbsoluteHttps(relative_url))
}

func init() {
	HtmlFuncBoot.Register(func(w *Web) {
		w.HtmlFunc["url"] = func(name string, a ...interface{}) string {
			return URLRev.Print(name, a...)
		}
	})
}
