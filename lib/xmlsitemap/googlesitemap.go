package xmlsitemap

import (
	"encoding/xml"
	"time"
)

const _xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Url     []url    `xml:"url"`
}

type url struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

func (ur UrlSet) Render() string {
	urls := []url{}
	for _, a_url := range ur {
		urls = append(urls, url{
			Loc:        a_url.Loc,
			LastMod:    a_url.LastMod.UTC().Format(time.RFC3339),
			ChangeFreq: a_url.ChangeFreq,
			Priority:   a_url.Priority,
		})
	}

	url_set := urlset{
		Xmlns: _xmlns,
		Url:   urls,
	}

	xmlbyte, _ := xml.MarshalIndent(url_set, "", "\t")
	return `<?xml version="1.0" encoding="UTF-8"?>` + string(xmlbyte)
}
