package logger

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"os"
	"time"
)

type LoggerInterface interface {
	Log(*webby.Web) string
}

type LoggerFunc func(*webby.Web) string

func (lo LoggerFunc) Log(w *webby.Web) string {
	return lo(w)
}

func LoggerStd(w *webby.Web) string {
	return fmt.Sprintf("%s %s %s %s ?%s IP:%s %v", w.Req.Proto, w.Req.Method,
		w.Req.Host, w.Req.URL.Path,
		w.Req.URL.RawQuery, w.RemoteAddr(), time.Now())
}

var DefaultLogger LoggerInterface = LoggerFunc(LoggerStd)

type StorageInterface interface {
	Store(string)
}

type StorageFunc func(string)

func (st StorageFunc) Store(str string) {
	st(str)
}

func DevNull(str string) {
	// Do nothing
}

var DefaultStorage StorageInterface = StorageFunc(DevNull)

const LogName = "Log"

// File Storage
type FileStorage struct {
	Path string
	Size int64
}

func (fi FileStorage) Store(str string) {
	if fi.Size == 0 || fi.Path == "" {
		return
	}

	// Init variables
	var file *os.File
	name := fi.Path + "/" + LogName

	// Check file size
	stat, err := os.Stat(name)
	if err == nil {
		if stat.Size() > fi.Size {
			curtime := time.Now()
			os.Rename(name, fmt.Sprintf("%s_%d_%d.txt", name, curtime.Unix(), curtime.UnixNano()))
		} else {
			file, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND, 0666)
		}
	}

	if file == nil {
		file, err = os.Create(name)
	}

	if err != nil {
		panic(err)
	}

	file.WriteString(str + "\r\n")

	// Never forget to close file.
	file.Close()
}

func init() {
	webby.PostBoot.Register(func(w *webby.Web) {
		DefaultStorage.Store(DefaultLogger.Log(w))
	})
}
