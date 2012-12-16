package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

func TestFormInputText(t *testing.T) {
	fmt.Println("InputText Test:\r\n")

	form := New(nil,
		&InputText{
			Name:       "text",
			MinChar:    1,
			MaxChar:    8,
			RegExpRule: "^([a-zA-Z]*)$",
			RegExpErr:  "Letters Only",
		},
		&InputText{
			Name:         "textmatch",
			MustMatch:    "text",
			MustMatchErr: "Does not match field above!",
		},
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"text":      []string{"hello"},
				"textmatch": []string{"hello"},
			},
			MultipartForm: &multipart.Form{},
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
				"text":      []string{"hello"},
				"textmatch": []string{"hellofail"},
			},
			MultipartForm: &multipart.Form{},
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
				"text":      []string{"hellohello"},
				"textmatch": []string{"hellohello"},
			},
			MultipartForm: &multipart.Form{},
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
				"text":      []string{"1234"},
				"textmatch": []string{"1234"},
			},
			MultipartForm: &multipart.Form{},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
