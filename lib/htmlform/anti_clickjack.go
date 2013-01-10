package htmlform

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

var (
	AntiClickJackExpire = 1 * time.Hour
)

type antiClickJack struct {
	key    string
	expire time.Time
}

var _antiClickJack = &antiClickJack{}

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

func setclickjack() {
	curtime := time.Now()
	_antiClickJack.key = fmt.Sprintf("%x%x", uint64ToByte(uint64(curtime.Unix())),
		uint64ToByte(uint64(curtime.UnixNano())))
	_antiClickJack.expire = curtime.Add(AntiClickJackExpire)
}

func init() {
	setclickjack()
	gob.Register(&inputClickJack{})
}

type inputClickJack struct {
	Value string
	error error
	lang  Lang
}

func (fo *inputClickJack) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="hidden" name="_anti-clickjack" class="anticlickjack" value="{{.Value}}"/>`
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *inputClickJack) Validate(values Values, files FileHeaders, single bool) error {
	fo.Value = values.Get("_anti-clickjack")

	if time.Now().Unix() > _antiClickJack.expire.Unix() {
		setclickjack()
	}

	if fo.Value != _antiClickJack.key {
		fo.Value = _antiClickJack.key
		return FormError(fo.lang["ErrAntiClickjack"])
	}

	return nil
}

func (fo *inputClickJack) GetName() string {
	return "_anti_clickjack"
}

func (fo *inputClickJack) SetError(err error) {
	fo.error = err
}

func (fo *inputClickJack) GetError() error {
	return fo.error
}

func (fo *inputClickJack) GetStruct() FormHandler {
	return fo
}

func (fo *inputClickJack) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *inputClickJack) GetLang() Lang {
	return fo.lang
}
