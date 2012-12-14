// Google XML Sitemap generator! Based on http://www.sitemaps.org/protocol.html
package xmlsitemap

import (
	"time"
)

type UrlSet []Url

type Url struct {
	Loc        string
	LastMod    time.Time
	ChangeFreq string
	Priority   float64
}
