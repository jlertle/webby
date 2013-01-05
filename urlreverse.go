package webby

import (
	"fmt"
)

type URLReverseMap map[string]string

type URLReverse struct {
	urls URLReverseMap
}

func (u *URLReverse) register(name, format string) {
	u.urls[name] = format
}

func (u *URLReverse) initUrls() {
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
	return fmt.Sprintf(u.urls[name], a...)
}

var URLRev = &URLReverse{}

func urlBootstrap(w *Web) {
	w.HtmlFunc["url"] = func(name string, a ...interface{}) string {
		return URLRev.Print(name, a...)
	}
}

func (w *Web) URLReverse(name string, a ...interface{}) string {
	return URLRev.Print(name, a...)
}

func init() {
	Boot.Register(urlBootstrap)
}
