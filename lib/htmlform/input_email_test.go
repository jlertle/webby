package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestInputEmail(t *testing.T) {
	fmt.Println("InputEmail Test:\r\n")

	form := New(nil,
		&InputEmail{
			Name: "email",
		},
		&InputEmail{
			Name:         "emailconfirm",
			MustMatch:    "email",
			MustMatchErr: "Must match field above!",
		},
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"_anti-CSRF":   []string{GetAntiCSRFKey()},
				"email":        []string{"himself@cj-jackson.com"},
				"emailconfirm": []string{"himself@cj-jackson.com"},
			},
		},
	}

	if !form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()

	web = &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"_anti-CSRF":   []string{GetAntiCSRFKey()},
				"email":        []string{"himself@cj-jackson.com"},
				"emailconfirm": []string{"himselfa@cj-jackson.com"},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()

	web = &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"_anti-CSRF":   []string{GetAntiCSRFKey()},
				"email":        []string{"himself_cj-jackson.com"},
				"emailconfirm": []string{"himself_cj-jackson.com"},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
