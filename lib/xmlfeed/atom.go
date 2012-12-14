package xmlfeed

import (
	"encoding/xml"
	"time"
)

const (
	_xmlns = "http://www.w3.org/2005/Atom"
)

// Atom Feed Structure
type atomFeed struct {
	XMLName  xml.Name    `xml:"feed"`
	Xmlns    string      `xml:"xmlns,attr"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Link     atomLink    `xml:"link"`
	Id       string      `xml:"id"`
	Updated  string      `xml:"updated"`
	Entry    []atomEntry `xml:"entry"`
}

func (atom atomFeed) render() string {
	output, _ := xml.MarshalIndent(atom, "", "  ")
	return "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>" + string(output)
}

// Atom Entry Structure
type atomEntry struct {
	Title   string   `xml:"title"`
	Link    atomLink `xml:"link"`
	Id      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Summary string   `xml:"summary"`
	Name    string   `xml:"author>name"`
	Email   string   `xml:"author>email"`
}

// Atom Link Structure
type atomLink struct {
	Link string `xml:"href,attr"`
}

// Generate Atom Feed!
func (channel Channel) Atom() string {
	atom_entry := []atomEntry{}
	for _, entry := range channel.Item {
		atom_entry = append(atom_entry, atomEntry{
			Title:   entry.Title,
			Link:    atomLink{entry.Link},
			Id:      entry.Link,
			Updated: entry.Updated.UTC().Format(time.RFC3339),
			Summary: entry.Description,
			Name:    entry.Name,
			Email:   entry.Email,
		})
	}

	atom_feed := atomFeed{
		Xmlns:    _xmlns,
		Title:    channel.Title,
		Subtitle: channel.Description,
		Link:     atomLink{channel.Link},
		Id:       channel.Link,
		Updated:  channel.Updated.UTC().Format(time.RFC3339),
		Entry:    atom_entry,
	}

	return atom_feed.render()
}
