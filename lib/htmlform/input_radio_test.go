package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

func TestInputRadio(t *testing.T) {
	fmt.Println("InputRadio Test:\r\n")

	form := New(nil,
		&InputRadio{
			Name:  "radio",
			Value: "car",
		},
		&InputRadio{
			Name:  "radio",
			Value: "motorbike",
		},
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"radio": []string{"car"},
			},
			MultipartForm: &multipart.Form{},
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
			MultipartForm: &multipart.Form{},
		},
	}

	form.IsValid(web)

	fmt.Println(form.Render())
	fmt.Println()
}
