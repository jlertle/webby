package webby

import (
	"compress/flate"
	"compress/gzip"
	"strings"
)

// Init Compression Buffer (Call Before Writing to Client)
func (w *Web) InitCompression() {
	if w.Req.Method == "HEAD" {
		return
	}

	if w.Env.Get("Connection") == "Upgrade" {
		return
	}

	for _, encoding := range strings.Split(w.Env.Get("Accept-Encoding"), ",") {
		encoding = strings.TrimSpace(strings.ToLower(encoding))
		switch encoding {
		case "gzip":
			w.pri.reswrite = gzip.NewWriter(w.web)
			w.Header().Set("Content-Encoding", encoding)
			return
		case "deflate":
			w.pri.reswrite, _ = flate.NewWriter(w.web, flate.DefaultCompression)
			w.Header().Set("Content-Encoding", encoding)
			return
		}
	}
}

func (w *Web) closeCompression() {
	switch t := w.pri.reswrite.(type) {
	case *gzip.Writer:
		t.Close()
	case *flate.Writer:
		t.Close()
	}
}
