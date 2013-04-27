package htmlform

import (
	"bytes"
	"encoding/gob"
)

type gobEval interface {
	Eval() FormHandlerExt
}

type gobFormHandler struct {
	Form interface{}
	Err  error
}

type gobForm struct {
	Fields   []gobFormHandler
	Lang     Lang
	AllowGet bool
}

func init() {
	gob.Register(gobForm{})
	gob.Register(gobFormHandler{})
}

func (f *Form) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	defer buf.Reset()

	form := gobForm{Lang: f.lang, AllowGet: f.allowGet}

	for _, value := range f.fields {
		value.SetLang(nil)
		form.Fields = append(form.Fields, gobFormHandler{value.GetStruct(), value.GetError()})
	}

	enc := gob.NewEncoder(buf)
	err := enc.Encode(form)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *Form) GobDecode(b []byte) error {
	buf := &bytes.Buffer{}
	defer buf.Reset()

	buf.Write(b)

	dec := gob.NewDecoder(buf)

	form := gobForm{}

	err := dec.Decode(&form)
	if err != nil {
		return err
	}

	f.lang = form.Lang
	f.allowGet = form.AllowGet

	for _, value := range form.Fields {
		avalue := value.Form.(gobEval).Eval()
		avalue.SetLang(form.Lang)
		avalue.SetError(value.Err)
		f.fields = append(f.fields, avalue)
	}

	return nil
}
