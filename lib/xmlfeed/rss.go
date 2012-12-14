package xmlfeed

import (
	"encoding/xml"
	"time"
)

// RSS Channel Structure
type rssChannel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Description   string    `xml:"description"`
	Link          string    `xml:"link"`
	LastBuildDate string    `xml:"lastBuildDate"`
	PubDate       string    `xml:"pubDate"`
	Ttl           int64     `xml:"ttl"`
	Item          []rssItem `xml:"item"`
}

func (rsschan rssChannel) render() string {
	output, _ := xml.MarshalIndent(rsschan, "", "  ")
	return "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n<rss version=\"2.0\">" + string(output) + "\n</rss>"
}

// RSS Item Structure
type rssItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

// Generate RSS2 Feed!
func (channel Channel) RSS() string {
	rss_item := []rssItem{}
	for _, item := range channel.Item {
		rss_item = append(rss_item, rssItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			Guid:        item.Link,
			PubDate:     item.PubDate.UTC().Format(time.RFC1123Z),
		})
	}

	rss := rssChannel{
		Title:         channel.Title,
		Description:   channel.Description,
		Link:          channel.Link,
		LastBuildDate: channel.LastBuildDate.UTC().Format(time.RFC1123Z),
		PubDate:       channel.PubDate.UTC().Format(time.RFC1123Z),
		Ttl:           channel.Ttl,
		Item:          rss_item,
	}
	return rss.render()
}
