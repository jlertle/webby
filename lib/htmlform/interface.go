package htmlform

import (
	"bytes"
)

type FormHandler interface {
	Render(buf *bytes.Buffer)
	Validate(values Values, files FileHeaders, single bool) error
	GetName() string
	SetError(err error)
	GetError() error
	GetStruct() FormHandler
	SetLang(lang Lang)
	GetLang() Lang
}
