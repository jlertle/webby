package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestSelect(t *testing.T) {
	fmt.Println("Select Test:\r\n")

	form := New(nil,
		&Select{
			Name:      "select",
			Mandatory: true,
			Options: []*Option{
				&Option{
					Name:  "Motorcycle",
					Value: "motorcycle",
				},
				&Option{
					Name:  "Car",
					Value: "car",
				},
			},
		},
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"_anti-CSRF": []string{GetAntiCSRFKey()},
				"select":     []string{"car"},
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
				"_anti-CSRF": []string{GetAntiCSRFKey()},
				"select":     []string{"motorcycle"},
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
				"_anti-CSRF": []string{GetAntiCSRFKey()},
				"select":     []string{""},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
