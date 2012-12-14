package xmlsitemap

import (
	"fmt"
	"testing"
	"time"
)

func TestSiteMap(t *testing.T) {
	urls := UrlSet{}

	urls = append(urls, Url{
		Loc:        "http://cj-jackson.com",
		LastMod:    time.Now(),
		ChangeFreq: "daily",
		Priority:   0.5,
	})
	fmt.Println(urls.Render())
}
