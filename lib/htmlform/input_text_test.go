package htmlform

import (
	"fmt"
	"github.com/cj-jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestFormInputText(t *testing.T) {
	fmt.Println("InputText Test:\r\n")

	form := New(
		NewInputText("text").Mandatory().MaxChar(8).RegExp("^([a-zA-Z]*)$", "Letters Only").Get(),
		NewInputText("textmatch").MustMatch("text", "Does not match field above!").Get(),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"text":      []string{"hello", "123"},
				"textmatch": []string{"hello", "hello"},
			},
		},
	}

	if !form.IsValid(web) {
		t.Fail()
	}

	if form.IsValidSlot(web, 1) {
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
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
