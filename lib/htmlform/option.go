package htmlform

import (
	"bytes"
	"encoding/gob"
)

type Option struct {
	Name     string
	Value    string
	Label    string
	Selected bool
}

func init() {
	gob.Register(Option{})
}

func (op *Option) Render(buf *bytes.Buffer) {
	const htmlstr = `<option {{if .IsValue}}value="{{.Value}}" {{end}}{{if .IsLabel}}label="{{.Label}}" {{end}}{{if .Selected}} selected {{end}}{{if .IsName}}>{{.Name}}</option>{{else}}/>{{end}}`
	htmlRender(buf, htmlstr, op)
}

func (op *Option) IsName() bool {
	return len(op.Name) > 0
}

func (op *Option) IsValue() bool {
	return len(op.Value) > 0
}

func (op *Option) IsLabel() bool {
	return len(op.Label) > 0
}
