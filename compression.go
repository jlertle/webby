package webby

import (
	"compress/flate"
	"compress/gzip"
	"strings"
)

// Init Compression Buffer (Call Before Writing to Client)
func (w *Web) InitCompression() {
	if w.Env.Get("Connection") == "Upgrade" {
		return
	}

	if w.Env.Get("Accept-Encoding") == "" {
		return
	}

	for _, encoding := range strings.Split(w.Env.Get("Accept-Encoding"), ",") {
		encoding = strings.TrimSpace(strings.ToLower(encoding))

		switch encoding {
		case "gzip":
			w.reswrite = gzip.NewWriter(w.webInterface)
			w.Header().Set("Content-Encoding", encoding)
			return
		case "deflate":
			w.reswrite, _ = flate.NewWriter(w.webInterface, flate.DefaultCompression)
			w.Header().Set("Content-Encoding", encoding)
			return
		}

	}
}

func (w *Web) closeCompression() {
	switch t := w.reswrite.(type) {
	case *gzip.Writer:
		t.Close()
	case *flate.Writer:
		t.Close()
	}
}
