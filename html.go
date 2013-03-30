package webby

import (
	"bytes"
	html "html/template"
	"io"
	"io/ioutil"
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

func (w *Web) bufToWeb(buf *bytes.Buffer) {
	io.Copy(w, buf)
	buf.Reset()
}

func (w *Web) parseHtml(htmlstr string, value_map interface{}, buf io.Writer) {
	if buf == nil {
		// To prevent headers from being sent too early.
		buf = &bytes.Buffer{}
		defer w.bufToWeb(buf.(*bytes.Buffer))
	}
	t, err := html.New("html").Funcs(w.HtmlFunc).Parse(htmlstr)
	w.Check(err)
	err = t.Execute(buf, value_map)
	w.Check(err)
}

// Parse HTML
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
func (w *Web) ParseHtml(htmlstr string, value_map interface{}) string {
	buf := &bytes.Buffer{}
	defer buf.Reset()
	w.parseHtml(htmlstr, value_map, buf)
	return buf.String()
}

// Parse HTML and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
func (w *Web) ParseHtmlSend(htmlstr string, value_map interface{}) {
	w.parseHtml(htmlstr, value_map, nil)
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
func (w *Web) GetHtmlFile(htmlfile string) string {
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

// Parse HTML File
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (w *Web) ParseHtmlFile(htmlfile string, value_map interface{}) string {
	return w.ParseHtml(w.GetHtmlFile(htmlfile), value_map)
}

// Parse HTML File and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (w *Web) ParseHtmlFileSend(htmlfile string, value_map interface{}) {
	w.ParseHtmlSend(w.GetHtmlFile(htmlfile), value_map)
}
