package webby

import (
	"testing"
)

func TestBootstrap(t *testing.T) {
	expected := 5
	count := 0

	bootItem := func(w *Web) {
		count++
	}

	bootCutItem := func(w *Web) {
		count++
		w.pri.cut = true
	}

	w := &Web{
		pri: &webPrivate{
			cut: false,
		},
	}

	boot := NewBootstrapReg(
		bootItem,
		bootItem,
		bootItem,
		bootItem,
		bootItem,
	)

	boot.Load(w)

	if count != expected {
		t.Fail()
	}

	expected = 3
	count = 0

	boot = NewBootstrapReg(
		bootItem,
		bootItem,
		bootCutItem,
		bootItem,
		bootItem,
	)

	boot.Load(w)

	if count != expected {
		t.Fail()
	}
}
