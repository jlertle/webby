// XML Feed Helper
package xmlfeed

import (
	"encoding/gob"
	"github.com/cj-jackson/webby"
	"time"
)

// Channel Structure
type Channel struct {
	Title         string
	Description   string
	Link          string
	LastBuildDate time.Time
	PubDate       time.Time
	Updated       time.Time
	Ttl           int64
	Item          []Item
}

// Item Structure
type Item struct {
	Title       string
	Description string
	Link        string
	PubDate     time.Time
	Updated     time.Time
	Name        string
	Email       string
}

func init() {
	gob.Register(Channel{})
	gob.Register(Item{})
}

type XmlFeed interface {
	Feed(*webby.Web) Channel
}

type AtomRouteHandler struct {
	XmlFeed
}

func (at AtomRouteHandler) View(w *webby.Web) {
	w.InitCompression()
	w.Print(at.Feed(w).Atom())
}

type RssRouteHandler struct {
	XmlFeed
}

func (rss RssRouteHandler) View(w *webby.Web) {
	w.InitCompression()
	w.Print(rss.Feed(w).RSS())
}
