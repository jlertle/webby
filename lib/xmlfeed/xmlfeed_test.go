package xmlfeed

import (
	"fmt"
	"testing"
	"time"
)

func TestXML(t *testing.T) {
	feed := Channel{
		Title:         "test",
		Description:   "test",
		Link:          "http://cj-jackson.com/",
		LastBuildDate: time.Now(),
		PubDate:       time.Now(),
		Updated:       time.Now(),
		Ttl:           1800,
		Item: []Item{
			Item{
				Title:       "test",
				Description: "test",
				Link:        "http://cj-jackson.com/test/",
				PubDate:     time.Now(),
				Updated:     time.Now(),
				Name:        "Christopher John Jackson",
				Email:       "himself@cj-jackson.com",
			},
		},
	}

	fmt.Print(feed.RSS(), "\r\n\r\n")
	fmt.Print(feed.Atom(), "\r\n\r\n")
}
