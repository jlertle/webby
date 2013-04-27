package htmlform

import (
	"bytes"
	"encoding/gob"
)

type InputRadio struct {
	Name     string
	Value    string
	Id       string
	Class    string
	Selected bool
}

func init() {
	gob.Register(InputRadio{})
}

func (fo *InputRadio) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="radio" name="{{.Name}}"
	value="{{.Value}}"
	{{if .IsId}}id="{{.Id}}"
	{{end}}{{if .IsClass}}class="{{.Class}}"
	{{end}}{{if .Selected}}checked
	{{end}}/>`
	htmlRender(buf, htmlstr, fo)
}

func (fo *InputRadio) Validate(val *Validation) error {
	values, _, _ := val.GetAll()
	if !values.Exist(fo.Name) {
		return nil
	}

	fo.Selected = values.Get(fo.Name) == fo.Value

	return nil
}

func (fo *InputRadio) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *InputRadio) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *InputRadio) GetName() string {
	return fo.Name
}

func (fo *InputRadio) SetError(err error) {
	// Do nothing
}

func (fo *InputRadio) GetError() error {
	return nil
}

func (fo *InputRadio) GetStruct() FormHandlerExt {
	return fo
}

func (fo *InputRadio) SetLang(lang Lang) {
	// Do nothing!
}

func (fo *InputRadio) GetLang() Lang {
	return nil
}

func (fo InputRadio) Eval() FormHandlerExt {
	return &fo
}
