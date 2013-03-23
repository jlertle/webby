package htmlform

import (
	"fmt"
	"testing"
)

func TestInputFile(t *testing.T) {
	fmt.Println("InputFile Test:\r\n")

	form := New(
		NewInputFile("file").Mime("image/jpeg", "image/png").Get(),
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
