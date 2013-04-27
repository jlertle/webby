package htmlform

import (
	"bytes"
)

type FormHandler interface {
	GetStruct() FormHandlerExt
}

type FormHandlerExt interface {
	Render(*bytes.Buffer)
	Validate(*Validation) error
	GetName() string
	SetError(err error)
	GetError() error
	GetStruct() FormHandlerExt
	SetLang(lang Lang)
	GetLang() Lang
}
