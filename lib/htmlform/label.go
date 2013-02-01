package htmlform

import (
	"bytes"
	"encoding/gob"
)

type Label struct {
	Name  string
	For   string
	Id    string
	Class string
}

func init() {
	gob.Register(Label{})
}

func (fo Label) Render(buf *bytes.Buffer) {
	const htmlstr = `{{if .IsFor}}<label for="{{.For}}"{{else}}<div{{end}}{{if .IsId}} id="{{.Id}}"{{end}}{{if .IsClass}} class="{{.Class}}"{{end}}>{{.Name}}{{if .IsFor}}</label>{{else}}</div>{{end}}`
	htmlRender(buf, htmlstr, fo)
}

func (fo Label) Validate(values Values, files FileHeaders, single bool) error {
	return nil
}

func (fo Label) GetName() string {
	return fo.Name
}

func (fo Label) SetError(err error) {
	// Do nothing
}

func (fo Label) GetError() error {
	return nil
}

func (fo Label) GetStruct() FormHandler {
	return fo
}

func (fo Label) IsFor() bool {
	return len(fo.For) > 0
}

func (fo Label) IsId() bool {
	return len(fo.Id) > 0
}

func (fo Label) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo Label) SetLang(lang Lang) {
	// Do nothing
}

func (fo Label) GetLang() Lang {
	return nil
}

func (fo Label) Eval() FormHandler {
	return fo
}
