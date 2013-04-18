package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestFormInputCheckbox(t *testing.T) {
	fmt.Println("InputCheckbox Test:\r\n")

	form := New(
		NewInputCheckbox("checkbox").Value("car").Get(),
		NewInputCheckbox("checkbox").Value("motorcycle").Mandatory().Get(),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"checkbox": []string{"car", "motorcycle"},
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
				"checkbox": []string{"car"},
			},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
