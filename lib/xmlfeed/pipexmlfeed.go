package xmlfeed

import (
	"time"
)

type PipeChannel struct {
	ch Channel
}

func NewChannel(title string) PipeChannel {
	return PipeChannel{
		Channel{Title: title, Ttl: 1800},
	}
}

func (ch PipeChannel) Description(description string) PipeChannel {
	ch.ch.Description = description
	return ch
}

func (ch PipeChannel) Link(link string) PipeChannel {
	ch.ch.Link = link
	return ch
}

func (ch PipeChannel) LastBuildDate(lastBuildDate time.Time) PipeChannel {
	ch.ch.LastBuildDate = lastBuildDate
	return ch
}

func (ch PipeChannel) PubDate(pubDate time.Time) PipeChannel {
	ch.ch.PubDate = pubDate
	return ch
}

func (ch PipeChannel) Updated(updated time.Time) PipeChannel {
	ch.ch.Updated = updated
	return ch
}

func (ch PipeChannel) Ttl(ttl int64) PipeChannel {
	ch.ch.Ttl = ttl
	return ch
}

func (ch PipeChannel) Item(item ...Item) PipeChannel {
	ch.ch.Item = item
	return ch
}

func (ch PipeChannel) AddItem(item Item) PipeChannel {
	ch.ch.Item = append(ch.ch.Item, item)
	return ch
}

func (ch PipeChannel) Get() Channel {
	return ch.ch
}

type PipeItem struct {
	it Item
}

func NewItem(title string) PipeItem {
	return PipeItem{
		Item{Title: title},
	}
}

func (it PipeItem) Description(description string) PipeItem {
	it.it.Description = description
	return it
}

func (it PipeItem) Link(link string) PipeItem {
	it.it.Link = link
	return it
}

func (it PipeItem) PubDate(pubDate time.Time) PipeItem {
	it.it.PubDate = pubDate
	return it
}

func (it PipeItem) Updated(updated time.Time) PipeItem {
	it.it.Updated = updated
	return it
}

func (it PipeItem) Name(name string) PipeItem {
	it.it.Name = name
	return it
}

func (it PipeItem) Email(email string) PipeItem {
	it.it.Email = email
	return it
}

func (it PipeItem) Get() Item {
	return it.it
}
