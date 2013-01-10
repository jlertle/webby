package htmlform

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

var (
	AntiCSRFExpire = 1 * time.Hour
)

type antiCSRF struct {
	key    string
	expire time.Time
}

var _antiCSRF = &antiCSRF{}

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

func setAntiCSRF() {
	curtime := time.Now()
	_antiCSRF.key = fmt.Sprintf("%x%x", uint64ToByte(uint64(curtime.Unix())),
		uint64ToByte(uint64(curtime.UnixNano())))
	_antiCSRF.expire = curtime.Add(AntiCSRFExpire)
}

func init() {
	setAntiCSRF()
	gob.Register(&inputCSRF{})
}

type inputCSRF struct {
	Value string
	error error
	lang  Lang
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

	if time.Now().Unix() > _antiCSRF.expire.Unix() {
		setAntiCSRF()
	}

	if fo.Value != _antiCSRF.key {
		fo.Value = _antiCSRF.key
		return FormError(fo.lang["ErrAntiCSRF"])
	}

	return nil
}

func (fo *inputCSRF) GetName() string {
	return "_anti_CSRF"
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
