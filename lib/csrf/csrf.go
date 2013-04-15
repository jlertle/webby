// provide Protection Against CSRF
package csrf

import (
	"github.com/cj-jackson/webby"
	html "html/template"
	"net/http"
)

const cookieName = "csrf_token"

var (
	// Just to avoid ugly url patterns!  The also include HEAD request!
	// Also GET and HEAD are both considered safe according to rfc2616, section 9.1.1 (http://tools.ietf.org/html/rfc2616.html#section-9.1.1)
	IncludeGetRequest = false
	// Modulisation is pretty useful for large site! Or when you want to specify the correct placement for csrf checking!
	Modulised = false
)

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

func init() {
	webby.HtmlFuncBoot.Register(func(w *webby.Web) {
		// Get CSRF Token input field
		w.HtmlFunc["csrf_token"] = func() html.HTML {
			const htmlstr = `<input type="hidden" name="{{.Name}}" class="{{.Name}}" value="{{.Value}}" />`
			return html.HTML(w.Html().Parse(htmlstr, getCookie(w)))
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

		Check{}.Boot(w)
	})
}
