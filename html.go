package webby

import (
	"bytes"
	html "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"
)

func init() {
	HtmlFuncBoot.Register(func(w *Web) {
		// HTML Marksafe
		w.HtmlFunc["html"] = func(str string) html.HTML {
			return html.HTML(str)
		}

		// HTML Attr MarkSafe
		w.HtmlFunc["htmlattr"] = func(str string) html.HTMLAttr {
			return html.HTMLAttr(str)
		}

		// JS Marksafe
		w.HtmlFunc["js"] = func(str string) html.JS {
			return html.JS(str)
		}

		// JS String Marksafe
		w.HtmlFunc["jsstr"] = func(str string) html.JSStr {
			return html.JSStr(str)
		}

		// CSS Marksafe
		w.HtmlFunc["css"] = func(str string) html.CSS {
			return html.CSS(str)
		}
	})
}

type Html struct {
	w *Web
}

func (w *Web) Html() Html {
	return Html{w}
}

func (h Html) render(htmlstr string, value_map interface{}, buf io.Writer) {
	w := h.w
	if buf == nil {
		// To prevent headers from being sent too early.
		buf = &bytes.Buffer{}
		defer func(buf *bytes.Buffer) {
			io.Copy(w, buf)
			buf.Reset()
		}(buf.(*bytes.Buffer))
	}
	t := html.Must(html.New("html").Funcs(w.HtmlFunc).Parse(htmlstr))
	err := t.Execute(buf, value_map)
	w.Check(err)
}

// Render HTML
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
func (h Html) Render(htmlstr string, value_map interface{}) string {
	buf := &bytes.Buffer{}
	defer buf.Reset()
	h.render(htmlstr, value_map, buf)
	return buf.String()
}

// Render HTML and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
func (h Html) RenderSend(htmlstr string, value_map interface{}) {
	h.render(htmlstr, value_map, nil)
}

type htmlFileCacheStruct struct {
	content string
	expire  time.Time
}

var htmlFileCache = struct {
	sync.Mutex
	m map[string]interface{}
}{m: map[string]interface{}{}}

var HtmlTemplateCacheExpire = 24 * time.Hour

// Get HTML File
//
// Note: Can also be used to get other kind of files. DO NOT USE THIS WITH LARGE FILES.
func (h Html) GetFile(htmlfile string) string {
	var content string
	var content_in_byte []byte
	var err error

	htmlFileCache.Lock()
	defer htmlFileCache.Unlock()

	switch t := htmlFileCache.m[htmlfile].(type) {
	case htmlFileCacheStruct:
		if time.Now().Unix() > t.expire.Unix() {
			goto getfile_and_cache
		}
		content = t.content
		goto return_content
	}

getfile_and_cache:
	content_in_byte, err = ioutil.ReadFile(htmlfile)
	if err != nil {
		return err.Error()
	}
	content = string(content_in_byte)
	if !DEBUG {
		htmlFileCache.m[htmlfile] = htmlFileCacheStruct{content, time.Now().Add(HtmlTemplateCacheExpire)}
	}

return_content:
	return content
}

// Render HTML File
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (h Html) RenderFile(htmlfile string, value_map interface{}) string {
	return h.Render(h.GetFile(htmlfile), value_map)
}

// Render HTML File and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (h Html) RenderFileSend(htmlfile string, value_map interface{}) {
	h.RenderSend(h.GetFile(htmlfile), value_map)
}

func (h Html) ParseFiles(filenames ...string) *html.Template {
	t := html.New("html").Funcs(h.w.HtmlFunc)

	for _, filename := range filenames {
		html.Must(t.Parse(h.GetFile(filename)))
	}

	return t
}

var htmlGlobLocker = struct {
	sync.Mutex
	filenames map[string][]string
}{
	filenames: map[string][]string{},
}

func (h Html) ParseGlob(pattern string) *html.Template {
	htmlGlobLocker.Lock()
	defer htmlGlobLocker.Unlock()

	if len(htmlGlobLocker.filenames[pattern]) > 0 {
		return h.ParseFiles(htmlGlobLocker.filenames[pattern]...)
	}

	filenames, err := filepath.Glob(pattern)
	h.w.Check(err)

	if !DEBUG {
		htmlGlobLocker.filenames[pattern] = filenames
	}

	return h.ParseFiles(filenames...)
}

type htmlDefault struct {
	filenames []string
	pattern   string
	template  *html.Template
}

func (h Html) init() {
	if h.w.pri.html == nil {
		h.w.pri.html = &htmlDefault{}
	}
}

func (h Html) SetDefaultFiles(filenames ...string) {
	h.init()
	h.w.pri.html.filenames = filenames
}

func (h Html) SetDefaultGlob(pattern string) {
	h.init()
	h.w.pri.html.pattern = pattern
}

func (h Html) Default() *html.Template {
	if h.w.pri.html == nil {
		panic(ErrorStr("HTML: Default Template is not set!"))
	}

	if h.w.pri.html.template != nil {
		goto gotoreturn
	}

	if len(h.w.pri.html.filenames) > 0 {
		h.w.pri.html.template = h.ParseFiles(h.w.pri.html.filenames...)
		goto gotoreturn
	}

	if h.w.pri.html.pattern != "" {
		h.w.pri.html.template = h.ParseGlob(h.w.pri.html.pattern)
		goto gotoreturn
	}

gotoreturn:
	return h.w.pri.html.template
}

func (h Html) DefaultRender(name string, data interface{}) string {
	b := &bytes.Buffer{}
	defer b.Reset()
	h.Default().ExecuteTemplate(b, name, data)
	return b.String()
}

func (h Html) DefaultRenderSend(name string, data interface{}) {
	b := &bytes.Buffer{}
	defer func(buf *bytes.Buffer) {
		io.Copy(h.w, buf)
		buf.Reset()
	}(b)
	h.Default().ExecuteTemplate(b, name, data)
}
