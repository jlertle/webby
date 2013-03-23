package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestFormInputPassword(t *testing.T) {
	fmt.Println("InputPassword Test:\r\n")

	form := New(
		NewInputPassword("password").Mandatory().MaxChar(8).RegExp("^([a-zA-Z]*)$", "Letters Only").Get(),
		NewInputPassword("passwordmatch").MustMatch("password", "Does not match field above!").Get(),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"password":      []string{"hello"},
				"passwordmatch": []string{"hello"},
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
				"password":      []string{"hello"},
				"passwordmatch": []string{"hellofail"},
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
				"password":      []string{"hellohello"},
				"passwordmatch": []string{"hellohello"},
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
				"password":      []string{"1234"},
				"passwordmatch": []string{"1234"},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
