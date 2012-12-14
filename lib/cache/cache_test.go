package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	SetAdv("test", "test", time.Now().Add(5*time.Second))

	switch Get("test").(type) {
	case nil:
		t.Fail()
	}

	time.Sleep(10 * time.Second)

	switch Get("test").(type) {
	case string:
		t.Fail()
	}
}
