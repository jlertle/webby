// HTML Form Helper!
package htmlform

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/CJ-Jackson/webby"
	"github.com/CJ-Jackson/webby/lib/htmlform/lang"
	html "html/template"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strconv"
)

type Lang map[string]string

func init() {
	gob.Register(Lang{})
}

var defaultLang = "EnglishGB"

func DefaultLang(str string) {
	defaultLang = str
}

var ErrorTemplate = `<p>{{.}}</p>`

type Values url.Values

func (va Values) Get(name string) string {
	if len(va[name]) <= 0 {
		return ""
	}
	return va[name][0]
}

func (va Values) Exist(name string) bool {
	if len(va[name]) <= 0 {
		return false
	}
	return true
}

type FileHeaders map[string][]*multipart.FileHeader

func (fi FileHeaders) Exist(name string) bool {
	if fi == nil {
		return false
	}

	if len(fi[name]) <= 0 {
		return false
	}

	return true
}

func (fi FileHeaders) GetSize(name string) int64 {
	size := int64(0)

	if !fi.Exist(name) {
		return size
	}

	file, err := fi[name][0].Open()
	if err != nil {
		return size
	}

	defer file.Close()

	// Not suppose to use seek to get filesize! But it worked perfectly!
	size, err = file.Seek(0, 2)
	if err != nil {
		return 0
	} else {
		file.Seek(0, 0)
	}

	return size
}

func (fi FileHeaders) GetContentType(name string) string {
	if !fi.Exist(name) {
		return ""
	}

	return fi[name][0].Header.Get("Content-Type")
}

type FormError string

func init() {
	gob.Register(FormError(""))
}

func (f FormError) Error() string {
	return string(f)
}

type Form struct {
	lang     Lang
	fields   []FormHandlerExt
	allowGet bool
}

func init() {
	gob.Register(Form{})
}

// Construct New Form Helper
func New(formhandlers ...FormHandler) *Form {
	f := &Form{fields: []FormHandlerExt{}, allowGet: false}

	for _, formhandler := range formhandlers {
		switch t := formhandler.(type) {
		case FormHandlerExt:
			f.fields = append(f.fields, t)
		default:
			f.fields = append(f.fields, formhandler.GetStruct())
		}
	}

	return f.Lang(defaultLang)
}

// Set Language
func (f *Form) Lang(lng string) *Form {
	f.lang = Lang(lang.Langs.Get(lng))

	for _, field := range f.fields {
		field.SetLang(f.lang)
	}

	return f
}

// Allow Get Request (Beware: Get Request can bypass CSRF protection.)
func (f *Form) AllowGet() *Form {
	f.allowGet = true
	return f
}

func (f *Form) Render() string {
	buf := &bytes.Buffer{}
	defer buf.Reset()
	for _, field := range f.fields {
		field.Render(buf)
	}
	return buf.String()
}

func (f *Form) RenderSlices() []string {
	buf := &bytes.Buffer{}
	var slices []string
	for _, field := range f.fields {
		field.Render(buf)
		slices = append(slices, buf.String())
		buf.Reset()
	}
	return slices
}

func (f *Form) isValid(w *webby.Web, values Values, files FileHeaders) bool {
	if f.allowGet {
		goto validation
	}

	if w.Req.Method == "GET" {
		return false
	}

validation:

	valid := true

	val := &Validation{values, files, false, CurVal("")}

	for _, field := range f.fields {
		field.SetError(nil)
		err := field.Validate(val)
		if err != nil {
			field.SetError(err)
			valid = false
		}
		val.CurVal = CurVal("")
	}

	return valid
}

func (f *Form) IsValid(w *webby.Web) bool {
	var files FileHeaders
	form := w.Form()
	values := Values(form.Value)
	files = nil
	if form.File != nil {
		files = FileHeaders(form.File)
	}
	return f.isValid(w, values, files)
}

// For the more complexed form!
func (f *Form) IsValidSlot(w *webby.Web, slot int) bool {
	var files FileHeaders
	form := w.Form().Slot(slot)
	values := Values(form.Value)
	files = nil
	if form.File != nil {
		files = FileHeaders(form.File)
	}
	return f.isValid(w, values, files)
}

// Very Useful for AJAX Validater that require server-side validation
func (f *Form) ValidateSingle(name, value, mime string) error {
	values := Values{
		name: []string{value},
	}
	mimeheader := textproto.MIMEHeader{}
	mimeheader.Add("Content-Type", mime)
	files := FileHeaders{
		name: []*multipart.FileHeader{&multipart.FileHeader{
			Header: mimeheader}},
	}

	val := &Validation{values, files, true, CurVal("")}

	for _, field := range f.fields {
		switch t := field.(type) {
		case Label:
			continue
		default:
			if t.GetName() == name {
				return t.Validate(val)
			}
		}
	}
	return FormError(f.lang["ErrFieldDoesNotExist"])
}

func (f *Form) getmap() map[string]FormHandler {
	themap := map[string]FormHandler{}
	for _, field := range f.fields {
		switch t := field.(type) {
		case Label:
			continue
		default:
			themap[t.GetName()] = t.GetStruct()
		}
	}
	return themap
}

// Useful for pure client-side validation (JavaScript)
func (f *Form) JSON() string {
	b, _ := json.Marshal(f.getmap())
	return string(b)
}

func htmlRender(buf *bytes.Buffer, htmlstr string, value_map interface{}) {
	t, err := html.New("html").Parse(htmlstr)
	if err != nil {
		buf.WriteString(err.Error())
	}
	err = t.Execute(buf, value_map)
	if err != nil {
		buf.WriteString(err.Error())
	}
}

// Convert String to int64
func toInt(number string) (int64, error) {
	return strconv.ParseInt(number, 10, 64)
}

func init() {
	webby.HtmlFuncBoot.Register(func(w *webby.Web) {
		// Render Form
		w.HtmlFunc["render_form"] = func(f *Form) string {
			return f.Render()
		}

		// Render Form Slices
		w.HtmlFunc["render_form_slices"] = func(f *Form) []string {
			return f.RenderSlices()
		}
	})
}

type CurVal string

func (cur CurVal) String() string {
	return string(cur)
}

func (cur CurVal) Int64() int64 {
	num, err := toInt(cur.String())
	if err != nil {
		return int64(0)
	}
	return num
}

func (cur CurVal) Int() int {
	return int(cur.Int64())
}

func (cur CurVal) Int32() int32 {
	return int32(cur.Int64())
}

func (cur CurVal) Int16() int16 {
	return int16(cur.Int64())
}

func (cur CurVal) Int8() int8 {
	return int8(cur.Int64())
}

func (cur CurVal) Uint64() uint64 {
	num, err := strconv.ParseUint(cur.String(), 10, 64)
	if err != nil {
		return uint64(0)
	}
	return num
}

func (cur CurVal) Uint() uint {
	return uint(cur.Uint64())
}

func (cur CurVal) Uint32() uint32 {
	return uint32(cur.Uint64())
}

func (cur CurVal) Uint16() uint16 {
	return uint16(cur.Uint64())
}

func (cur CurVal) Uint8() uint8 {
	return uint8(cur.Uint64())
}

func (cur CurVal) Float64() float64 {
	num, err := strconv.ParseFloat(cur.String(), 64)
	if err != nil {
		return float64(0)
	}
	return num
}

func (cur CurVal) Float32() float32 {
	return float32(cur.Float64())
}

type Validation struct {
	Val    Values
	Files  FileHeaders
	Single bool
	CurVal CurVal
}

func (val *Validation) GetAll() (Values, FileHeaders, bool) {
	return val.Val, val.Files, val.Single
}

type ExtraFunc func(*Validation) error
