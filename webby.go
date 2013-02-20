// Simple Web Application Server Framework!
package webby

import (
	"bufio"
	"fmt"
	html "html/template"
	"io"
	"net"
	"net/http"
	"runtime/debug"
)

// Debug Mode
var DEBUG = false

var RootView RouteHandler = BootRoute{Boot, Route}

type webInterface interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
}

// The Framework Structure
type Web struct {
	webInterface
	// Error Code
	Status int
	// Server Environment Variables
	Env http.Header
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
		webInterface: res.(webInterface),
		Status:       http.StatusOK,
		Env:          req.Header,
		Req:          req,
		Meta:         map[string]interface{}{},
		Param:        Param{},
		HtmlFunc:     html.FuncMap{},
		Session:      nil,
		path:         req.URL.Path,
		curpath:      "",
		cut:          false,
		firstWrite:   true,
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack()
			if DEBUG {
				web.Status = 500
				web.Println("500 Internal Server Error")
				web.Printf("\r\n%s\r\n\r\n%s", r, debug.Stack())
				return
			}
			web.Error500()
		}
	}()

	web.reswrite = web.webInterface

	web.initTrueHost()
	web.initTrueRemoteAddr()
	web.initSession()
	web.Header().Set("Content-Encoding", "plain")

	defer web.closeCompression()

	web.debugStart()
	defer web.debugEnd()

	MainBoot.Load(web)

	if web.CutOut() {
		return
	}

	RootView.View(web)

	if web.CutOut() {
		return
	}

	Error500(web)
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
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

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
// Note: Use Status properly to set error code! As this disable compression!
func (web *Web) WriteHeader(num int) {
	web.cut = true

	if web.firstWrite {
		web.firstWrite = false
	}

	web.webInterface.WriteHeader(num)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the HTTP server library
// will not do anything else with the connection.
// It becomes the caller's responsibility to manage
// and close the connection.
func (web *Web) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	web.cut = true
	return web.webInterface.Hijack()
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

func (web *Web) debuginfo(a string) {
	if !DEBUG {
		return
	}
	fmt.Printf("--\r\n %s  %s, %s, %s, %s, ?%s \r\n--\r\n",
		a, web.Req.Proto, web.Req.Method,
		web.Req.Host, web.Req.URL.Path,
		web.Req.URL.RawQuery)
}

func (web *Web) debugStart() {
	web.debuginfo("START")
}

func (web *Web) debugEnd() {
	web.debuginfo("END  ")
}
