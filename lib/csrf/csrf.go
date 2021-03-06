package csrf

import (
	"github.com/CJ-Jackson/webby"
	html "html/template"
	"net/http"
)

const cookieName = "csrf_token"

// Just to avoid ugly url patterns!  The also include HEAD request!
// Also GET and HEAD are both considered safe according to rfc2616,
// section 9.1.1 (http://tools.ietf.org/html/rfc2616.html#section-9.1.1)
var IncludeGetRequest = false

// Modulisation is pretty useful for large site! Or when you want to
// specify the correct placement for csrf checking!
var Modulised = false

func getCookie(w *webby.Web) *http.Cookie {
	cookie, err := w.Cookie(cookieName).Get()
	if err != nil {
		cookie, _ = w.Cookie(cookieName).Value(webby.KeyGen()).Month().SaveRes().SaveReq().Get()
	}
	return cookie
}

func formCheck(w *webby.Web) {
	form := w.Form()
	if len(form.Value) <= 0 {
		return
	}

	if len(form.Value[cookieName]) <= 0 {
		w.Error403()
		return
	}

	if form.Value[cookieName][0] != getCookie(w).Value {
		w.Error403()
		return
	}
}

/*
CSRF Bootstrap Handler
*/
type Check struct{}

func (_ Check) Boot(w *webby.Web) {
	if IncludeGetRequest {
		goto form_check
	}

	switch w.Req.Method {
	case "GET", "HEAD":
		return
	}

form_check:

	formCheck(w)
}

var _check = Check{}

func init() {
	webby.HtmlFuncBoot.Register(func(w *webby.Web) {
		// Get CSRF Token input field
		w.HtmlFunc["csrf_token"] = func() html.HTML {
			const htmlstr = `<input type="hidden" name="{{.Name}}" class="{{.Name}}" value="{{.Value}}" />`
			return html.HTML(w.Html().Render(htmlstr, getCookie(w)))
		}

		// Get CSRF Token Key
		w.HtmlFunc["csrf_token_key_only"] = func() string {
			return getCookie(w).Value
		}
	})

	webby.MainBoot.Register(func(w *webby.Web) {
		if Modulised {
			return
		}

		_check.Boot(w)
	})
}
