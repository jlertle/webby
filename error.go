package webby

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"
)

func printPanic(buf io.Writer, w *Web, r interface{}, stack []byte) {
	fmt.Fprintf(buf, "\r\n%s, %s, %s, %s, ?%s IP:%s\r\n",
		w.Req.Proto, w.Req.Method,
		w.Req.Host, w.Req.URL.Path,
		w.Req.URL.RawQuery, w.Req.RemoteAddr)

	fmt.Fprintf(buf, "\r\n%s\r\n\r\n%s", r, stack)

	fmt.Fprintln(buf, "\r\nRequest Header:")
	fmt.Fprintln(buf, w.Req.Header)

	w.ParseForm()

	fmt.Fprintln(buf, "\r\nForm Values:")
	fmt.Fprintln(buf, w.Req.Form)

	fmt.Fprintln(buf, "\r\nForm Values (Multipart):")
	fmt.Fprintln(buf, w.Req.MultipartForm)

	fmt.Fprintln(buf, "\r\nTime:")
	fmt.Fprintln(buf, time.Now())
}

// Check for Error
func (w *Web) Check(err error) {
	Check(err)
}

// Check for Error
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

type PanicHandler interface {
	Panic(*Web, interface{}, []byte)
}

type PanicConsole struct{}

func (_ PanicConsole) Panic(w *Web, r interface{}, stack []byte) {
	ErrPrint(r, "\r\n", string(stack))
}

const panicFileExt = ".txt"

type PanicFile struct {
	Path string
}

func (p PanicFile) Panic(w *Web, r interface{}, stack []byte) {
	filename := p.Path + fmt.Sprintf("/%d_%d", time.Now().Unix(), time.Now().UnixNano()) + panicFileExt
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	printPanic(file, w, r, stack)
}

var DefaultPanicHandler PanicHandler = PanicConsole{}

type Errors struct {
	E403 func(w *Web)
	E404 func(w *Web)
	E500 func(w *Web)
}

// Overridable Error403 Function
//
// Note:  Overriding is useful for custom 403 page
var Error403 = func(w *Web) {
	w.Print("<h1>403 Forbidden</h1>")
}

func (w *Web) Error403() {
	w.Status = 403
	w.Errors.E403(w)
}

// Overridable Error404 Function
//
// Note:  Overriding is useful for custom 404 page
var Error404 = func(w *Web) {
	w.Print("<h1>404 Not Found</h1>")
}

func (w *Web) Error404() {
	w.Status = 404
	w.Errors.E404(w)
}

// Overridable Error500 Function
//
// Note:  Overriding is useful for custom 500 page
var Error500 = func(w *Web) {
	w.Print("<h1>500 Internal Server Error</h1>")
}

func (w *Web) Error500() {
	w.Status = 500
	w.Errors.E500(w)
}

type ErrorStr string

func (e ErrorStr) Error() string {
	return "Error: " + string(e)
}

func ErrPrint(a ...interface{}) {
	fmt.Fprint(os.Stderr, a...)
}

func ErrPrintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func ErrPrintln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func (w *Web) recover() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		DefaultPanicHandler.Panic(w, r, stack)
		if DEBUG {
			w.Status = 500
			w.Println("500 Internal Server Error")
			printPanic(w, w, r, stack)
			return
		}
		w.Error500()
	}
}
