package webby

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"
)

func printPanic(buf io.Writer, w *Web, r interface{}, stack []byte) {
	printF := func(format string, a ...interface{}) {
		fmt.Fprintf(buf, format, a...)
	}

	printLn := func(a ...interface{}) {
		fmt.Fprintln(buf, a...)
	}

	printF("\r\n%s, %s, %s, %s, ?%s IP:%s\r\n",
		w.Req.Proto, w.Req.Method,
		w.Req.Host, w.Req.URL.Path,
		w.Req.URL.RawQuery, w.Req.RemoteAddr)

	printF("\r\n%s\r\n\r\n%s", r, stack)

	printLn("\r\nRequest Header:")
	printLn(w.Req.Header)

	w.ParseForm()

	printLn("\r\nForm Values:")
	printLn(w.Req.Form)

	printLn("\r\nForm Values (Multipart):")
	printLn(w.Req.MultipartForm)

	printLn("\r\nTime:")
	printLn(time.Now())
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

// Write error to Stderr.
type PanicConsole struct{}

func (_ PanicConsole) Panic(w *Web, r interface{}, stack []byte) {
	ErrPrint(r, "\r\n", string(stack))
}

const panicFileExt = ".txt"

// Write error to new file.
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

// Execute Error 403 (Forbidden)
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

// Execute Error 404 (Not Found)
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

// Execute Error 500 (Internal Server Error)
func (w *Web) Error500() {
	w.Status = 500
	w.Errors.E500(w)
}

// Custom String Data Type, Implement error interface.
type ErrorStr string

func (e ErrorStr) Error() string {
	return "Error: " + string(e)
}

// Print formats using the default formats for its operands and writes to standard error output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func ErrPrint(a ...interface{}) {
	fmt.Fprint(os.Stderr, a...)
}

// Printf formats according to a format specifier and writes to standard error output.
// It returns the number of bytes written and any write error encountered.
func ErrPrintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// Println formats using the default formats for its operands and writes to standard error output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
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
