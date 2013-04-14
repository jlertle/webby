package htmlform

import (
	"github.com/cj-jackson/webby"
)

type RouteInterface interface {
	Init(*webby.Web) bool
	FetchForm(*webby.Web) *Form
	Get(*webby.Web, *Form)
	PostPass(*webby.Web)
	PostFail(*webby.Web)
}

// Only suitable with the simplist of form!
type RouteHandler struct {
	RouteInterface
}

func (ro RouteHandler) View(w *webby.Web) {
	truth := ro.Init(w)

	if w.CutOut() {
		return
	}

	if !truth {
		w.Error403()
		return
	}

	var form *Form

	switch t := w.Session.(type) {
	case Form:
		form = &t
		w.DestroySession()
	case *Form:
		form = t
		w.DestroySession()
	default:
		form = ro.FetchForm(w)
	}

	if w.Req.Method == "POST" {
		goto post
	}

	if w.Req.URL.RawQuery != "" {
		if form.IsValid(w) {
			goto pass
		} else {
			goto get
		}
	}

get:
	ro.Get(w, form)
	return

post:
	if !form.IsValid(w) {
		goto fail
	}

pass:
	ro.PostPass(w)
	return

fail:
	w.SetSession(form)
	ro.PostFail(w)
}
