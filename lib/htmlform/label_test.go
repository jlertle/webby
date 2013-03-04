package htmlform

import (
	"fmt"
	"testing"
)

func TestLabelForm(t *testing.T) {
	fmt.Println("Label Test: \r\n")

	form := New(
		Label{
			Name:  "Label",
			For:   "label",
			Id:    "Label",
			Class: "Label",
		},
	)

	fmt.Println(form.Render())

	form = New(
		Label{
			Name: "Label",
			For:  "label",
		},
	)

	fmt.Println(form.Render())

	form = New(
		Label{
			Name: "Label",
		},
	)

	fmt.Println(form.Render())

	fmt.Println()
}
