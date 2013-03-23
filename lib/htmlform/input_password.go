package htmlform

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"strings"
)

type InputPassword struct {
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
	error        error
	lang         Lang
}

func init() {
	gob.Register(InputPassword{})
}

func (fo *InputPassword) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="password" name="{{.Name}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}{{if .IsValue}}value="{{.Value}}" {{end}}{{if .IsMaxChar}}maxlength="{{.MaxChar}}" {{end}}{{if .IsRegExp}}pattern="{{.RegExpRule}}" {{end}}/>`
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *InputPassword) Validate(val *Validation) error {
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

	return nil
}

func (fo *InputPassword) IsValue() bool {
	return len(fo.Value) > 0
}

func (fo *InputPassword) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *InputPassword) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *InputPassword) IsMaxChar() bool {
	return fo.MaxChar > 0
}

func (fo *InputPassword) IsRegExp() bool {
	return len(fo.RegExpRule) > 0
}

func (fo *InputPassword) GetName() string {
	return fo.Name
}

func (fo *InputPassword) SetError(err error) {
	fo.error = err
}

func (fo *InputPassword) GetError() error {
	return fo.error
}

func (fo *InputPassword) GetStruct() FormHandler {
	return fo
}

func (fo *InputPassword) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *InputPassword) GetLang() Lang {
	return fo.lang
}

func (fo InputPassword) Eval() FormHandler {
	return &fo
}
