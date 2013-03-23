package htmlform

import (
	"bytes"
	"encoding/gob"
	"regexp"
	"strings"
)

// From http://www.w3.org/TR/html5/states-of-the-type-attribute.html#valid-e-mail-address
var email_rule = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")

type InputEmail struct {
	Name         string
	Value        string
	Id           string
	Class        string
	MustMatch    string
	MustMatchErr string
	error        error
	lang         Lang
}

func init() {
	gob.Register(InputEmail{})
}

func (fo *InputEmail) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="email" name="{{.Name}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}{{if .IsValue}}value="{{.Value}}" {{end}}/>`
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo InputEmail) Validate(val *Validation) error {
	values, _, _ := val.GetAll()
	if !values.Exist(fo.Name) {
		return FormError(fo.lang["ErrMandatory"])
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

	if !email_rule.MatchString(fo.Value) {
		return FormError(fo.lang["ErrInvalidEmailAddress"])
	}

	return nil
}

func (fo *InputEmail) IsValue() bool {
	return len(fo.Value) > 0
}

func (fo *InputEmail) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *InputEmail) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *InputEmail) GetName() string {
	return fo.Name
}

func (fo *InputEmail) SetError(err error) {
	fo.error = err
}

func (fo *InputEmail) GetError() error {
	return fo.error
}

func (fo *InputEmail) GetStruct() FormHandler {
	return fo
}

func (fo *InputEmail) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *InputEmail) GetLang() Lang {
	return fo.lang
}

func (fo InputEmail) Eval() FormHandler {
	return &fo
}
