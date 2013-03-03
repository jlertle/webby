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

	form := New(nil,
		&InputPassword{
			Name:       "password",
			MinChar:    1,
			MaxChar:    8,
			RegExpRule: "^([a-zA-Z]*)$",
			RegExpErr:  "Letters Only",
		},
		&InputPassword{
			Name:         "passwordmatch",
			MustMatch:    "password",
			MustMatchErr: "Does not match field above!",
		},
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
