// XML Feed Helper
package xmlfeed

import (
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
