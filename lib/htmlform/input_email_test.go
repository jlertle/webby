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

	form := New(
		NewInputEmail("email").Get(),
		NewInputEmail("emailconfirm").MustMatch("email", "Must match field above!").Get(),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
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
