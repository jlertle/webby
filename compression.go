package webby

import (
	"compress/flate"
	"compress/gzip"
	"strings"
)

// Compression Enabled
var CompressionEnabled = true

func (web *Web) initCompression() {
	web.Header().Set("Content-Encoding", "plain")

	if !CompressionEnabled {
		return
	}

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

// Kill Compression
func (web *Web) KillCompression() {
	web.Header().Set("Content-Encoding", "plain")
	web.reswrite = web.Res
	web.firstWrite = false
	web.cut = true
}
