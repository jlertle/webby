package htmlform

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/CJ-Jackson/webby"
	"github.com/CJ-Jackson/webby/lib/cache"
	"time"
)

var (
	AntiCSRFExpire         = 1 * time.Hour
	AntiCSRFJavaScriptMode = false
)

type antiCSRF string

// Convert Unsigned 64-bit Int to Bytes.
func uint64ToByte(num uint64) [8]byte {
	var buf [8]byte
	buf[0] = byte(num >> 0)
	buf[1] = byte(num >> 8)
	buf[2] = byte(num >> 16)
	buf[3] = byte(num >> 24)
	buf[4] = byte(num >> 32)
	buf[5] = byte(num >> 40)
	buf[6] = byte(num >> 48)
	buf[7] = byte(num >> 56)
	return buf
}

func getAntiCSRF() string {
	switch t := cache.Get("_antiCsrf").(type) {
	case antiCSRF:
		return string(t)
	}

	curtime := time.Now()
	key := antiCSRF(fmt.Sprintf("%x%x", uint64ToByte(uint64(curtime.Unix())),
		uint64ToByte(uint64(curtime.UnixNano()))))
	cache.SetAdv("_antiCsrf", key, curtime.Add(AntiCSRFExpire))

	return string(key)
}

// Get AntiCSRF Key
func GetAntiCSRFKey() string {
	return getAntiCSRF()
}

type inputCSRF struct {
	Value string
	error error
	lang  Lang
}

func init() {
	gob.Register(inputCSRF{})
	gob.Register(antiCSRF(""))
}

func (fo *inputCSRF) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="hidden" name="_anti-CSRF" class="antiCSRF" value="{{.Value}}"/>`
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *inputCSRF) Validate(values Values, files FileHeaders, single bool) error {
	fo.Value = values.Get("_anti-CSRF")

	currentKey := getAntiCSRF()

	if fo.Value != currentKey {
		fo.Value = currentKey
		return FormError(fo.lang["ErrAntiCSRF"])
	}

	return nil
}

func (fo *inputCSRF) GetName() string {
	return "_anti-CSRF"
}

func (fo *inputCSRF) SetError(err error) {
	fo.error = err
}

func (fo *inputCSRF) GetError() error {
	return fo.error
}

func (fo *inputCSRF) GetStruct() FormHandler {
	return fo
}

func (fo *inputCSRF) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *inputCSRF) GetLang() Lang {
	return fo.lang
}

func (fo inputCSRF) Eval() FormHandler {
	return &fo
}

type CSRFRouteHandler struct {
	Key string `json:"key"`
}

func (_ CSRFRouteHandler) View(w *webby.Web) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	if !w.IsAjaxRequest() {
		w.Error404()
		return
	}
	csrf := CSRFRouteHandler{Key: getAntiCSRF()}
	enc := json.NewEncoder(w)
	enc.Encode(csrf)
}
