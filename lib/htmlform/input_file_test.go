package htmlform

import (
	"fmt"
	"testing"
)

func TestInputFile(t *testing.T) {
	fmt.Println("InputFile Test:\r\n")

	form := New(nil,
		&InputFile{
			Name: "file",
			Mime: []string{"image/jpeg", "image/png"},
		},
	)

	fmt.Println(form.Render())
	fmt.Println()

	if form.ValidateSingle("file", "", "image/jpeg") != nil {
		t.Fail()
	}

	if form.ValidateSingle("file", "", "image/bmp") == nil {
		t.Fail()
	}
}
