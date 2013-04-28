// Google XML Sitemap generator! Based on http://www.sitemaps.org/protocol.html
package xmlsitemap

import (
	"time"
)

type UrlSet []UrlInterface

type Url struct {
	Loc        string
	LastMod    time.Time
	ChangeFreq string
	Priority   float64
}

func (u Url) Get() Url {
	return u
}

type UrlInterface interface {
	Get() Url
}

// Chainable version of Url
type PipeUrl struct {
	url Url
}

func NewUrl() PipeUrl {
	return PipeUrl{Url{}}
}

func (pi PipeUrl) Get() Url {
	return pi.url
}

func (pi PipeUrl) Loc(loc string) PipeUrl {
	pi.url.Loc = loc
	return pi
}

func (pi PipeUrl) LastMod(lastMod time.Time) PipeUrl {
	pi.url.LastMod = lastMod
	return pi
}

func (pi PipeUrl) ChangeFreq(changeFreq string) PipeUrl {
	pi.url.ChangeFreq = changeFreq
	return pi
}

func (pi PipeUrl) Priority(priority float64) PipeUrl {
	pi.url.Priority = priority
	return pi
}
