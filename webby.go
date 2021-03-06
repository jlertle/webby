package webby

import (
	"bufio"
	html "html/template"
	"io"
	"net"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"
	"time"
)

// Debug Mode
var DEBUG = false

func init() {
	HtmlFuncBoot.Register(func(w *Web) {
		w.HtmlFunc["Debug"] = func() bool {
			return DEBUG
		}

		w.HtmlFunc["NotDebug"] = func() bool {
			return !DEBUG
		}
	})
}

var RootView RouteHandler = FuncToRouteHandler(func(w *Web) {
	appMiddlewares := AppMiddlewares.Init(w)
	appMiddlewares.Pre()
	defer appMiddlewares.Post()
	if w.CutOut() {
		return
	}

	Boot.Load(w)
	if w.CutOut() {
		return
	}

	Route.Load(w)
})

type web interface {
	http.ResponseWriter
}

type webPrivate struct {
	path       string
	curpath    string
	reswrite   io.Writer
	cut        bool
	firstWrite bool
	cmd        map[string]func(interface{}) interface{}
	html       *htmlDefault
	session    *SessionAdv
}

// The Framework Structure, it's implement the interfaces of 'net/http.ResponseWriter',
// 'net/http.Hijacker', 'net/http.Flusher' and 'net/http.Handler'
type Web struct {
	// Error Code
	Status int
	// Server Environment Variables
	Env http.Header
	// Request
	Req *http.Request
	// Meta, useful for storing login credentail
	Meta map[string]interface{}
	// Well same as meta, but for string data type only! Useful for storing user country code!
	Str map[string]string
	// Used by router for storing data of named group in RegExpRule
	Param Param
	// Function to load in html template system.
	HtmlFunc html.FuncMap
	// For holding session!
	Session interface{}
	// Errors
	Errors *Errors
	// Time Location
	TimeLoc *time.Location
	// Time Format
	TimeFormat string
	web
	pri *webPrivate
}

// HTTP Handler
func (_ Web) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	w := &Web{
		web:        res.(web),
		Status:     http.StatusOK,
		Env:        req.Header,
		Req:        req,
		Meta:       map[string]interface{}{},
		Str:        map[string]string{},
		Param:      Param{},
		HtmlFunc:   html.FuncMap{},
		Session:    nil,
		TimeLoc:    DefaultTimeLoc,
		TimeFormat: DefaultTimeFormat,
		Errors: &Errors{
			E403: Error403,
			E404: Error404,
			E500: Error500,
		},
		pri: &webPrivate{
			path:       req.URL.Path,
			curpath:    "",
			cut:        false,
			firstWrite: true,
			cmd:        map[string]func(interface{}) interface{}{},
		},
	}

	w.initWriter()
	w.initTrueHost()
	w.initTrueRemoteAddr()
	w.initTruePath()
	w.initSecure()
	w.initSession()

	defer w.recover()
	defer w.closeCompression()

	defer PostBoot.Load(w)

	w.debugStart()
	defer w.debugEnd()

	HtmlFuncBoot.Load(w)

	if w.CutOut() {
		return
	}

	mainMiddleware := MainMiddlewares.Init(w)
	mainMiddleware.Pre()
	defer mainMiddleware.Post()

	if w.CutOut() {
		return
	}

	MainBoot.Load(w)

	if w.CutOut() {
		return
	}

	w.RouteDealer(RootView)

	if w.CutOut() {
		return
	}

	w.Error500()
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (w *Web) Header() http.Header {
	return w.web.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (w *Web) Write(data []byte) (int, error) {
	w.pri.cut = true

	if w.pri.firstWrite {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(data))
		}

		w.pri.firstWrite = false
		w.WriteHeader(w.Status)
	}

	return w.pri.reswrite.Write(data)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
//
// Note: Use Status property to set error code! As this disable compression!
func (w *Web) WriteHeader(num int) {
	w.pri.cut = true

	if w.pri.firstWrite {
		w.pri.firstWrite = false
	}

	w.web.WriteHeader(num)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the HTTP server library
// will not do anything else with the connection.
// It becomes the caller's responsibility to manage
// and close the connection.
func (w *Web) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.pri.cut = true

	hj, ok := w.web.(http.Hijacker)
	if ok {
		return hj.Hijack()
	}

	return nil, nil, ErrorStr("Connection is not Hijackable")
}

// Flush sends any buffered data to the client.
func (w *Web) Flush() {
	fl, ok := w.web.(http.Flusher)
	if ok {
		fl.Flush()
	}
}

// true if output was sent to client, otherwise false!
func (w *Web) CutOut() bool {
	return w.pri.cut
}

func (w *Web) debuginfo(a string) {
	if !DEBUG {
		return
	}
	ErrPrintf("--\r\n %s  %s, %s, %s, %s, ?%s IP:%s \r\n--\r\n",
		a, w.Req.Proto, w.Req.Method,
		w.Req.Host, w.Req.URL.Path,
		w.Req.URL.RawQuery, w.Req.RemoteAddr)
}

func (w *Web) debugStart() {
	w.debuginfo("START")
}

func (w *Web) debugEnd() {
	w.debuginfo("END  ")
}

var _web = Web{}

func secure(res http.ResponseWriter, req *http.Request) {
	req.Header.Set("X-Secure-Mode", "1")
	_web.ServeHTTP(res, req)
}

func nonsecure(res http.ResponseWriter, req *http.Request) {
	req.Header.Del("X-Secure-Mode")
	_web.ServeHTTP(res, req)
}

// Start Http Server
func StartHttp(addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(nonsecure))
}

// Start Http Server with TLS
func StartHttpTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(secure))
}

// Start FastCGI Server
func StartFastCGI(l net.Listener) error {
	if l == nil {
		os.Stderr = nil
	}
	return fcgi.Serve(l, http.HandlerFunc(nonsecure))
}

// Start CGI, disables Stderr completely. (Due to the way how IIS handlers Stderr)
func StartCGI() error {
	os.Stderr = nil
	return cgi.Serve(http.HandlerFunc(nonsecure))
}
