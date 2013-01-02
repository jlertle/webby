// Simple Web Application Server Framework!
package webby

import (
	"fmt"
	html "html/template"
	"io"
	"net/http"
)

// Debug Mode
var DEBUG = false

// The Framework Structure
type Web struct {
	// Error Code
	Status int
	// Server Environment Variables
	Env http.Header
	// Responder, use Web.Print(), Web.Printf() or Web.Println() to output
	Res http.ResponseWriter
	// Request
	Req *http.Request
	// Meta, useful for storing login credentail
	Meta map[string]interface{}
	// Used by router for storing data of named group in RegExpRule
	Param Param
	// Function to load in html template system.
	HtmlFunc html.FuncMap
	// For holding session!
	Session    interface{}
	path       string
	curpath    string
	reswrite   io.Writer
	cut        bool
	firstWrite bool
}

// HTTP Handler
func (_ Web) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	web := &Web{
		Status:     http.StatusOK,
		Env:        req.Header,
		Res:        res,
		Req:        req,
		Meta:       map[string]interface{}{},
		Param:      Param{},
		HtmlFunc:   html.FuncMap{},
		Session:    nil,
		path:       req.URL.Path,
		curpath:    "",
		reswrite:   res,
		cut:        false,
		firstWrite: true,
	}

	web.initTrueHost()
	web.initTrueRemoteAddr()
	web.initSession()
	web.Header().Set("Content-Encoding", "plain")

	defer web.closeCompression()

	Boot.Load(web)

	if web.CutOut() {
		return
	}

	Route.Load(web)

	if web.CutOut() || web.IsWebSocketRequest() {
		return
	}

	Error500(web)
}

// HTTP Response Header
func (web *Web) Header() http.Header {
	return web.Res.Header()
}

// Write bytes to Client (http web server or browser)
func (web *Web) Write(data []byte) (int, error) {
	web.cut = true

	if web.firstWrite {
		if web.Header().Get("Content-Type") == "" {
			web.Header().Set("Content-Type", http.DetectContentType(data))
		}

		web.firstWrite = false
		web.WriteHeader(web.Status)
	}

	return web.reswrite.Write(data)
}

// Output Http Header, use Status properly to set error code! As this disable compression!
func (web *Web) WriteHeader(num int) {
	web.cut = true

	if web.firstWrite {
		web.firstWrite = false
	}

	web.Res.WriteHeader(num)
}

// Print formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (web *Web) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(web, a...)
}

// Printf formats according to a format specifier and writes to client (http web server or browser).
// It returns the number of bytes written and any write error encountered.
func (web *Web) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(web, format, a...)
}

// Println formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (web *Web) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(web, a...)
}

// true if output was sent to client, otherwise false!
func (web *Web) CutOut() bool {
	return web.cut
}

// Cut, useful for serving files!
func (web *Web) Cut() {
	web.cut = true
}
