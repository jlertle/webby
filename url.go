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

func (u Url) Reverse(name string, a ...interface{}) string {
	return URLRev.Print(name, a...)
}

func init() {
	HtmlFuncBoot.Register(func(w *Web) {
		w.HtmlFunc["url"] = func(name string, a ...interface{}) string {
			return URLRev.Print(name, a...)
		}
	})
}
