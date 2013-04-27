package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"net/http"
	"net/url"
	"testing"
)

func TestInputRadio(t *testing.T) {
	fmt.Println("InputRadio Test:\r\n")

	form := New(
		NewInputRadio("radio").Value("car"),
		NewInputRadio("radio").Value("motorbike"),
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"radio": []string{"car"},
			},
		},
	}

	form.IsValid(web)

	fmt.Println(form.Render())
	fmt.Println()

	web = &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"radio": []string{"motorbike"},
			},
		},
	}

	form.IsValid(web)

	fmt.Println(form.Render())
	fmt.Println()
}
