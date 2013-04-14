package htmlform

import (
	"fmt"
	"github.com/cj-jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestSelect(t *testing.T) {
	fmt.Println("Select Test:\r\n")

	form := New(
		NewSelect("select").Mandatory().Options(
			NewOption("Motorcycle").Value("motorcycle").Get(),
			NewOption("Car").Value("car").Get(),
		).Get(),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"select": []string{"car"},
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
				"select": []string{"motorcycle"},
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
				"select": []string{""},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
