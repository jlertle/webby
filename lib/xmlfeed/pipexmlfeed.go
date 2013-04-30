package xmlfeed

import (
	"time"
)

type PipeChannel struct {
	Ch Channel
}

func NewChannel(title string) *PipeChannel {
	return &PipeChannel{
		Channel{Title: title, Ttl: 1800},
	}
}

func (ch *PipeChannel) Description(description string) *PipeChannel {
	ch.Ch.Description = description
	return ch
}

func (ch *PipeChannel) Link(link string) *PipeChannel {
	ch.Ch.Link = link
	return ch
}

func (ch *PipeChannel) LastBuildDate(lastBuildDate time.Time) *PipeChannel {
	ch.Ch.LastBuildDate = lastBuildDate
	return ch
}

func (ch *PipeChannel) PubDate(pubDate time.Time) *PipeChannel {
	ch.Ch.PubDate = pubDate
	return ch
}

func (ch *PipeChannel) Updated(updated time.Time) *PipeChannel {
	ch.Ch.Updated = updated
	return ch
}

func (ch *PipeChannel) Ttl(ttl int64) *PipeChannel {
	ch.Ch.Ttl = ttl
	return ch
}

func (ch *PipeChannel) Item(items ...ItemInterface) *PipeChannel {
	if ch.Ch.Item == nil {
		ch.Ch.Item = []Item{}
	}
	for _, item := range items {
		ch.Ch.Item = append(ch.Ch.Item, item.Get())
	}
	return ch
}

func (ch *PipeChannel) AddItem(item ItemInterface) *PipeChannel {
	ch.Ch.Item = append(ch.Ch.Item, item.Get())
	return ch
}

func (ch *PipeChannel) Get() Channel {
	return ch.Ch
}

func (ch PipeChannel) Atom() string {
	return ch.Ch.Atom()
}

func (ch PipeChannel) RSS() string {
	return ch.Ch.RSS()
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
