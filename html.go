package webby

import (
	"bytes"
	html "html/template"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

func htmlraw(str string) html.HTML {
	return html.HTML(str)
}

func htmlattr(str string) html.HTMLAttr {
	return html.HTMLAttr(str)
}

func js(str string) html.JS {
	return html.JS(str)
}

func jsstr(str string) html.JSStr {
	return html.JSStr(str)
}

func html_bootstrap(w *Web) {
	w.HtmlFunc["html"] = htmlraw
	w.HtmlFunc["htmlattr"] = htmlattr
	w.HtmlFunc["js"] = js
	w.HtmlFunc["jsstr"] = jsstr
}

func init() {
	Boot.Register(html_bootstrap)
}

func (w *Web) parseHtml(htmlstr string, value_map interface{}, buf io.Writer) {
	if buf == nil {
		buf = w
	}
	t, err := html.New("html").Funcs(w.HtmlFunc).Parse(htmlstr)
	if err != nil {
		buf.Write([]byte(err.Error()))
	}
	err = t.Execute(buf, value_map)
	if err != nil {
		buf.Write([]byte(err.Error()))
	}
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
func (web *Web) ParseHtmlFile(htmlfile string, value_map interface{}) string {
	return web.ParseHtml(web.GetHtmlFile(htmlfile), value_map)
}

// Parse HTML File and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (web *Web) ParseHtmlFileSend(htmlfile string, value_map interface{}) {
	web.ParseHtmlSend(web.GetHtmlFile(htmlfile), value_map)
}
