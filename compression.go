package webby

import (
	"compress/flate"
	"compress/gzip"
	"strings"
)

// Init Compression Buffer
func (web *Web) InitCompression() {
	if web.Env.Get("Connection") == "Upgrade" {
		return
	}

	if web.Env.Get("Accept-Encoding") == "" {
		return
	}

	for _, encoding := range strings.Split(web.Env.Get("Accept-Encoding"), ",") {
		encoding = strings.TrimSpace(strings.ToLower(encoding))

		switch encoding {
		case "gzip":
			web.reswrite = gzip.NewWriter(web.Res)
			web.Header().Set("Content-Encoding", encoding)
			return
		case "deflate":
			web.reswrite, _ = flate.NewWriter(web.Res, flate.DefaultCompression)
			web.Header().Set("Content-Encoding", encoding)
			return
		}

	}
}

func (web *Web) closeCompression() {
	switch t := web.reswrite.(type) {
	case *gzip.Writer:
		t.Close()
	case *flate.Writer:
		t.Close()
	}
}
