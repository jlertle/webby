package urlutil

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	if Slug("Hello World") != "hello-world" {
		t.Fail()
	}
	if Slug("Hello World69") != "hello-world69" {
		t.Fail()
	}
	if Slug("Hello           World????????") != "hello-world-" {
		t.Fail()
	}
	if Slug("aáäâeéëeiíiîoóöoôuúüuunç·/_,:;") != "aaaaeeeeiiiiooooouuuuunc-" {
		t.Fail()
	}
}
