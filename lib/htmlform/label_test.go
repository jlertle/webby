package htmlform

import (
	"fmt"
	"testing"
)

func TestLabelForm(t *testing.T) {
	fmt.Println("Label Test: \r\n")

	form := New(
		NewLabel("Label").For("label").Id("Label").Class("Label"),
	)

	fmt.Println(form.Render())

	form = New(
		NewLabel("Label").For("label").Get(),
	)

	fmt.Println(form.Render())

	form = New(
		NewLabel("Label").Get(),
	)

	fmt.Println(form.Render())

	fmt.Println()
}
