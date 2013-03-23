package htmlform

import (
	"bytes"
)

type FormHandler interface {
	Render(*bytes.Buffer)
	Validate(*Validation) error
	GetName() string
	SetError(err error)
	GetError() error
	GetStruct() FormHandler
	SetLang(lang Lang)
	GetLang() Lang
}
