package htmlform

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"strings"
)

type InputText struct {
	Name         string
	Value        string
	Id           string
	Class        string
	MinChar      int
	MaxChar      int
	RegExpRule   string
	RegExpErr    string
	MustMatch    string
	MustMatchErr string
	extra        func(*Validation) error
	error        error
	lang         Lang
}

func init() {
	gob.Register(InputText{})
}

func (fo *InputText) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="text" name="{{.Name}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}{{if .IsValue}}value="{{.Value}}" {{end}}{{if .IsMaxChar}}maxlength="{{.MaxChar}}" {{end}}{{if .IsRegExp}}pattern="{{.RegExpRule}}" {{end}}/>`
	if fo.error != nil {
		htmlRender(buf, ErrorTemplate, fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *InputText) Validate(val *Validation) error {
	values, _, _ := val.GetAll()
	if !values.Exist(fo.Name) {
		if fo.MinChar >= 1 {
			return FormError(fo.lang["ErrMandatory"])
		}
		return nil
	}

	fo.Value = strings.TrimSpace(values.Get(fo.Name))

	var MatchValue string
	if len(fo.MustMatch) > 0 {
		if !values.Exist(fo.MustMatch) {
			return FormError(fo.lang["ErrMustMatchMissing"])
		}
		MatchValue = values.Get(fo.MustMatch)
	}

	if len(MatchValue) <= 0 {
		goto skipmatch
	}

	if MatchValue != fo.Value {
		return FormError(fo.MustMatchErr)
	}

skipmatch:

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

	var truth bool
	var err error

	if len(fo.RegExpRule) <= 0 {
		goto skiprule
	}

	truth, err = regexp.MatchString(fo.RegExpRule, fo.Value)

	if err != nil {
		return FormError(err.Error())
	}

	if !truth {
		return FormError(fo.RegExpErr)
	}

skiprule:

	if fo.extra == nil {
		goto skipextra
	}

	err = fo.extra(val)
	if err != nil {
		return err
	}

skipextra:

	return nil
}

func (fo *InputText) IsValue() bool {
	return len(fo.Value) > 0
}

func (fo *InputText) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *InputText) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *InputText) IsMaxChar() bool {
	return fo.MaxChar > 0
}

func (fo *InputText) IsRegExp() bool {
	return len(fo.RegExpRule) > 0
}

func (fo *InputText) GetName() string {
	return fo.Name
}

func (fo *InputText) SetError(err error) {
	fo.error = err
}

func (fo *InputText) GetError() error {
	return fo.error
}

func (fo *InputText) GetStruct() FormHandler {
	return fo
}

func (fo *InputText) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *InputText) GetLang() Lang {
	return fo.lang
}

func (fo InputText) Eval() FormHandler {
	return &fo
}
