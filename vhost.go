package webby

import (
	"strings"
)

// Use host name as string (e.g example.com)
type VHost map[string]BootRoute

func (v VHost) View(w *Web) {
	for host, bootroute := range v {
		if len(host) > len(w.Req.Host) {
			continue
		}
		if strings.ToLower(host) == strings.ToLower(w.Req.Host[:len(host)]) {
			bootroute.View(w)
			return
		}
	}

	w.Error404()
}
