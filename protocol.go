package webby

import (
	"strings"
)

// Demuxer for Protocol. Implement RouteHandler interface.
type Protocol struct {
	ALL, HTTP, HTTPS RouteHandler
}

func (pr Protocol) View(w *Web) {
	switch strings.ToLower(strings.Split(w.Req.Proto, "/")[0]) {
	case "http":
		if pr.HTTP != nil {
			w.RouteDealer(pr.HTTP)
			return
		}
	case "shttp", "https":
		if pr.HTTPS != nil {
			w.RouteDealer(pr.HTTPS)
			return
		}
	}

	if pr.ALL != nil {
		w.RouteDealer(pr.ALL)
		return
	}

	w.Error404()
	return
}

// Chainable version of Protocol
type PipeProtocol struct {
	pr Protocol
}

// PipeProtocol constructor
func NewProtocol() PipeProtocol {
	return PipeProtocol{Protocol{}}
}

// Get Protocol
func (pi PipeProtocol) Get() Protocol {
	return pi.pr
}

// Set Http
func (pi PipeProtocol) Http(http RouteHandler) PipeProtocol {
	pi.pr.HTTP = http
	return pi
}

// Set Https
func (pi PipeProtocol) Https(https RouteHandler) PipeProtocol {
	pi.pr.HTTPS = https
	return pi
}

// Set All
func (pi PipeProtocol) All(all RouteHandler) PipeProtocol {
	pi.pr.ALL = all
	return pi
}

// Alais of All
func (pi PipeProtocol) Any(any RouteHandler) PipeProtocol {
	return pi.All(any)
}

// Alais of All
func (pi PipeProtocol) Fallback(fallback RouteHandler) PipeProtocol {
	return pi.All(fallback)
}
