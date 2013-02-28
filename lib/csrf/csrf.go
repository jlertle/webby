// provide Protection Against CSRF
package csrf

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	html "html/template"
	"net/http"
	"time"
)

const cookieName = "csrf_token"

var IncludeGetRequest = false

// Convert Unsigned 64-bit Int to Bytes.
func uint64ToByte(num uint64) [8]byte {
	var buf [8]byte
	buf[0] = byte(num >> 0)
	buf[1] = byte(num >> 8)
	buf[2] = byte(num >> 16)
	buf[3] = byte(num >> 24)
	buf[4] = byte(num >> 32)
	buf[5] = byte(num >> 40)
	buf[6] = byte(num >> 48)
	buf[7] = byte(num >> 56)
	return buf
}

func genKey() string {
	curtime := time.Now()
	return fmt.Sprintf("%x%x", uint64ToByte(uint64(curtime.Unix())),
		uint64ToByte(uint64(curtime.UnixNano())))
}

func getCookie(w *webby.Web) *http.Cookie {
	cookie, err := w.Req.Cookie(cookieName)
	if err != nil {
		cookie = &http.Cookie{
			Name:    cookieName,
			Value:   genKey(),
			Expires: time.Now().AddDate(0, 1, 0),
		}
		w.SetCookie(cookie)
		w.Req.AddCookie(cookie)
	}
	return cookie
}

func fail(w *webby.Web) {
	w.WriteHeader(403)
	w.Println("403: Forbidden: Failed CSRF Validation!")
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

func init() {
	webby.MainBoot.Register(func(w *webby.Web) {
		w.HtmlFunc["csrf_token"] = func() html.HTML {
			const htmlstr = `<input type="hidden" name="{{.Name}}" class="{{.Name}}" value="{{.Value}}" />`
			return html.HTML(w.ParseHtml(htmlstr, getCookie(w)))
		}

		w.HtmlFunc["csrf_token_key_only"] = func() string {
			return getCookie(w).Value
		}

		if IncludeGetRequest {
			goto parse_form
		}

		if w.Req.Method == "GET" {
			return
		}

	parse_form:

		w.ParseForm()

		if w.Req.MultipartForm != nil {
			multipartCheck(w)
			return
		}

		formCheck(w)
	})
}
