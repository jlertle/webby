package htmlform

import (
	"bytes"
	"encoding/gob"
)

type Select struct {
	Name      string
	Id        string
	Class     string
	Options   []*Option
	Mandatory bool
	error     error
	lang      Lang
}

func init() {
	gob.Register(Select{})
}

func (fo *Select) Render(buf *bytes.Buffer) {
	const (
		htmlstr       = `<select name="{{.Name}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}>`
		htmlstr_close = "</select>"
	)
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
	for _, option := range fo.Options {
		option.Render(buf)
	}
	buf.WriteString(htmlstr_close)
}

func (fo *Select) Validate(val *Validation) error {
	values, _, _ := val.GetAll()
	if !values.Exist(fo.Name) {
		return FormError(fo.lang["ErrSelectValueMissing"])
	}

	value := values.Get(fo.Name)
	truth := false
	for _, option := range fo.Options {
		option.Selected = false
		if option.Value == value {
			option.Selected = true
			truth = true
		}
	}

	if truth {
		return nil
	}

	if !fo.Mandatory {
		return nil
	}

	return FormError(fo.lang["ErrSelectOptionIsRequired"])
}

func (fo *Select) GetName() string {
	return fo.Name
}

func (fo *Select) SetError(err error) {
	fo.error = err
}

func (fo *Select) GetError() error {
	return fo.error
}

func (fo *Select) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *Select) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *Select) GetStruct() FormHandler {
	return fo
}

func (fo *Select) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *Select) GetLang() Lang {
	return fo.lang
}

func (fo Select) Eval() FormHandler {
	return &fo
}
