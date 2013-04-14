package htmlform

import (
	"fmt"
	"github.com/cj-jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestTextarea(t *testing.T) {
	fmt.Println("Textarea Test:\r\n")

	form := New(
		NewTextarea("text").Mandatory().MaxChar(8).Get(),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"text": []string{"hello"},
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
				"text": []string{"hellohello"},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
