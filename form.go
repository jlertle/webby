package webby

import (
	"mime/multipart"
	"net/url"
)

type Form struct {
	Value url.Values
	File  map[string][]*multipart.FileHeader
}

func (f *Form) TrimPrefix(prefix string) *Form {
	form := &Form{Value: url.Values{}}
	prefixLen := len(prefix)
	for key, value := range f.Value {
		if len(key) <= prefixLen {
			continue
		}
		if prefix != key[:prefixLen] {
			continue
		}
		form.Value[key[prefixLen:]] = value
	}
	if f.File == nil {
		return form
	}
	form.File = map[string][]*multipart.FileHeader{}
	for key, value := range f.File {
		if len(key) <= prefixLen {
			continue
		}
		if prefix != key[:prefixLen] {
			continue
		}
		form.File[key[prefixLen:]] = value
	}
	return form
}

func (f *Form) Slot(slot int) *Form {
	form := &Form{Value: url.Values{}}
	for key, value := range f.Value {
		if len(value) > slot {
			form.Value[key] = append(form.Value[key], value[slot])
		}
	}
	if f.File == nil {
		return form
	}
	form.File = map[string][]*multipart.FileHeader{}
	for key, value := range f.File {
		if len(value) > slot {
			form.File[key] = append(form.File[key], value[slot])
		}
	}
	return form
}

func (w *Web) Form() *Form {
	w.ParseForm()
	form := &Form{}
	if w.Req.MultipartForm != nil {
		form.Value = url.Values(w.Req.MultipartForm.Value)
		form.File = w.Req.MultipartForm.File
	} else {
		form.Value = w.Req.Form
		form.File = nil
	}
	return form
}

func (w *Web) FormTrimPrefix(prefix string) *Form {
	return w.Form().TrimPrefix(prefix)
}

func (w *Web) FormSlot(slot int) *Form {
	return w.Form().Slot(slot)
}
