package htmlform

import (
	"bytes"
	"fmt"
	"strings"
)

type Textarea struct {
	Name    string
	Value   string
	Id      string
	Class   string
	MinChar int
	MaxChar int
	Rows    int
	Cols    int
	error   error
	lang    Lang
}

func (fo *Textarea) Render(buf *bytes.Buffer) {
	const htmlstr = `<textarea name="{{.Name}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}rows="{{.Rows}}" cols="{{.Cols}}" >{{.Value}}</textarea>`
	if fo.Rows <= 0 {
		fo.Rows = 4
	}
	if fo.Cols <= 0 {
		fo.Cols = 25
	}
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *Textarea) Validate(values Values, files FileHeaders, single bool) error {
	if !values.Exist(fo.Name) {
		if fo.MinChar >= 1 {
			return FormError(fo.lang["ErrMandatory"])
		}
		return nil
	}

	fo.Value = strings.TrimSpace(values.Get(fo.Name))

	if fo.MinChar <= 0 {
		goto skipmin
	}

	if len(fo.Value) < fo.MinChar {
		if fo.Value == "" {
			return FormError(fo.lang["ErrMandatory"])
		}
		if fo.MinChar == 1 || fo.Value[0] == ' ' ||
			fo.Value[0] == '\r' || fo.Value[0] == '\n' {
			return FormError(fo.lang["ErrMandatory"])
		} else {
			return FormError(fmt.Sprintf(fo.lang["ErrMinChar"], fo.MinChar))
		}
	}

skipmin:

	if fo.MaxChar <= 0 {
		goto skipmax
	}

	if len(fo.Value) > fo.MaxChar {
		return FormError(fmt.Sprintf(fo.lang["ErrMaxChar"], fo.MaxChar))
	}

skipmax:

	return nil
}

func (fo *Textarea) IsValue() bool {
	return len(fo.Value) > 0
}

func (fo *Textarea) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *Textarea) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *Textarea) GetName() string {
	return fo.Name
}

func (fo *Textarea) SetError(err error) {
	fo.error = err
}

func (fo *Textarea) GetError() error {
	return fo.error
}

func (fo *Textarea) GetStruct() FormHandler {
	return fo
}

func (fo *Textarea) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *Textarea) GetLang() Lang {
	return fo.lang
}
