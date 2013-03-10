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
			pr.HTTP.View(w)
			return
		}
	case "shttp", "https":
		if pr.HTTPS != nil {
			pr.HTTPS.View(w)
			return
		}
	}

	if pr.ALL != nil {
		pr.ALL.View(w)
		return
	}

	w.Error404()
	return
}

// Chainable version of Protocol
type PipeProtocol struct {
	pr Protocol
}

func NewProtocol() PipeProtocol {
	return PipeProtocol{Protocol{}}
}

func (pi PipeProtocol) Get() Protocol {
	return pi.pr
}

func (pi PipeProtocol) Http(http RouteHandler) PipeProtocol {
	pi.pr.HTTP = http
	return pi
}

func (pi PipeProtocol) Https(https RouteHandler) PipeProtocol {
	pi.pr.HTTPS = https
	return pi
}

func (pi PipeProtocol) All(all RouteHandler) PipeProtocol {
	pi.pr.ALL = all
	return pi
}

func (pi PipeProtocol) Any(any RouteHandler) PipeProtocol {
	return pi.All(any)
}

func (pi PipeProtocol) Fallback(fallback RouteHandler) PipeProtocol {
	return pi.All(fallback)
}
