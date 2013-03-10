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
