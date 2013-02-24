package webby

import (
	"fmt"
	"os"
	"time"
)

// Check for Error
func (web *Web) Check(err error) {
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

func (_ PanicConsole) Panic(web *Web, r interface{}, stack []byte) {
	fmt.Print(r, "\r\n", string(stack))
}

const panicFileExt = ".txt"

type PanicFile struct {
	Path string
}

func (p PanicFile) Panic(web *Web, r interface{}, stack []byte) {
	filename := p.Path + fmt.Sprintf("/%d_%d", time.Now().Unix(), time.Now().UnixNano()) + panicFileExt
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "\r\n%s, %s, %s, %s, ?%s IP:%s\r\n",
		web.Req.Proto, web.Req.Method,
		web.Req.Host, web.Req.URL.Path,
		web.Req.URL.RawQuery, web.Req.RemoteAddr)

	fmt.Fprintf(file, "\r\n%s\r\n\r\n%s", r, stack)

	fmt.Fprintln(file, "\r\nRequest Header:")
	fmt.Fprintln(file, web.Req.Header)

	web.ParseForm()

	fmt.Fprintln(file, "\r\nForm Values:")
	fmt.Fprintln(file, web.Req.Form)

	fmt.Fprintln(file, "\r\nForm Values (Multipart):")
	fmt.Fprintln(file, web.Req.MultipartForm)

	fmt.Fprintln(file, "\r\nTime:")
	fmt.Fprintln(file, time.Now())
}

var DefaultPanicHandler PanicHandler = PanicConsole{}
