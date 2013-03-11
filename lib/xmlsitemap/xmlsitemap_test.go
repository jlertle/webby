package xmlsitemap

import (
	"fmt"
	"testing"
	"time"
)

func TestSiteMap(t *testing.T) {
	urls := UrlSet{}

	NewUrl().Loc("http://cj-jackson.com").LastMod(
		time.Now()).ChangeFreq("daily").Priority(0.5).Get()

	urls = append(urls, NewUrl().Loc("http://cj-jackson.com").LastMod(
		time.Now()).ChangeFreq("daily").Priority(0.5).Get())

	fmt.Println(urls.Render())
}
