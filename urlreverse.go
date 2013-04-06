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

func (w *Web) URLReverse(name string, a ...interface{}) string {
	return URLRev.Print(name, a...)
}

func init() {
	HtmlFuncBoot.Register(func(w *Web) {
		w.HtmlFunc["url"] = func(name string, a ...interface{}) string {
			return URLRev.Print(name, a...)
		}
	})
}
