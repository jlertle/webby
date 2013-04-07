package htmlform

import (
	"github.com/CJ-Jackson/webby"
	"strings"
)

type Put struct {
	form     *Form
	url      string
	failFunc func(w *webby.Web)
	passFunc func(w *webby.Web)
}

func (f *Form) Put() Put {
	return Put{
		form: f,
	}
}

func (p Put) Url(url string) Put {
	p.url = url
	return p
}

func (p Put) Fail(fail func(w *webby.Web)) Put {
	p.failFunc = fail
	return p
}

func (p Put) Pass(pass func(w *webby.Web)) Put {
	p.passFunc = pass
	return p
}

func (p Put) Validate(w *webby.Web) {
	var data struct {
		Pass bool   `json:"pass"`
		Url  string `json:"url"`
		Html string `json:"html"`
	}
	w.Status = 201
	data.Pass = true
	data.Url = p.url

	if !p.form.IsValid(w) {
		w.Status = 202
		data.Pass = false
		data.Html = strings.Join(p.form.RenderSlices(), "<br />")
		if p.failFunc != nil {
			p.failFunc(w)
		}
	} else {
		if p.passFunc != nil {
			p.passFunc(w)
		}
	}

	w.JsonSend(data)
}
