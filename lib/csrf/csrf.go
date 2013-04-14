// provide Protection Against CSRF
package csrf

import (
	"github.com/cj-jackson/webby"
	html "html/template"
	"net/http"
	"time"
)

const cookieName = "csrf_token"

var (
	// Just to avoid ugly url patterns!  The also include HEAD request!
	// Also GET and HEAD are both considered safe according to rfc2616, section 9.1.1 (http://tools.ietf.org/html/rfc2616.html#section-9.1.1)
	IncludeGetRequest = false
	// Modulisation is pretty useful for large site! Or when you want to specify the correct placement for csrf checking!
	Modulised = false
)

func genKey() string {
	return webby.KeyGen()
}

func getCookie(w *webby.Web) *http.Cookie {
	cookie, err := w.GetCookie(cookieName)
	if err != nil {
		cookie = w.NewCookie(cookieName).Value(genKey()).Expires(time.Now().AddDate(0, 1, 0)).SaveRes().SaveReq().Get()
	}
	return cookie
}

func fail(w *webby.Web) {
	w.Error403()
}

func multipartCheck(w *webby.Web) {
	if len(w.Req.MultipartForm.Value[cookieName]) <= 0 {
		fail(w)
		return
	}

	if w.Req.MultipartForm.Value[cookieName][0] != getCookie(w).Value {
		fail(w)
		return
	}
}

func formCheck(w *webby.Web) {
	if len(w.Req.Form) <= 0 {
		return
	}

	if len(w.Req.Form[cookieName]) <= 0 {
		fail(w)
		return
	}

	if w.Req.Form[cookieName][0] != getCookie(w).Value {
		fail(w)
		return
	}
}

type Check struct{}

func (_ Check) Boot(w *webby.Web) {
	if IncludeGetRequest {
		goto parse_form
	}

	switch w.Req.Method {
	case "GET", "HEAD":
		return
	}

parse_form:

	w.ParseForm()

	if w.Req.MultipartForm != nil {
		multipartCheck(w)
		return
	}

	formCheck(w)
}

func init() {
	webby.HtmlFuncBoot.Register(func(w *webby.Web) {
		// Get CSRF Token input field
		w.HtmlFunc["csrf_token"] = func() html.HTML {
			const htmlstr = `<input type="hidden" name="{{.Name}}" class="{{.Name}}" value="{{.Value}}" />`
			return html.HTML(w.ParseHtml(htmlstr, getCookie(w)))
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
