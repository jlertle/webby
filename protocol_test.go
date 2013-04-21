package webby

import (
	"net/http"
	"testing"
)

type ProtocolDummy struct {
	Protocol
}

func (pr *ProtocolDummy) Http() {
	pr.W.Param.Set("protocol", "HTTP")
}

func (pr *ProtocolDummy) Https() {
	pr.W.Param.Set("protocol", "HTTPS")
}

func TestProtocol(t *testing.T) {
	w := &Web{
		Param: Param{},
		Env:   http.Header{},
		Req: &http.Request{
			Proto: "HTTP/1.1",
		},
		pri: &webPrivate{
			cut: false,
		},
	}

	protocol := func() string {
		return w.Param.Get("protocol")
	}

	w.RouteDealer(&ProtocolDummy{})

	if protocol() != "HTTP" {
		t.Fail()
	}

	w.Req.Proto = "SHTTP/1.3"

	w.RouteDealer(&ProtocolDummy{})

	if protocol() != "HTTPS" {
		t.Fail()
	}

	w.Req.Proto = "HTTPS/1.3"

	w.RouteDealer(&ProtocolDummy{})

	if protocol() != "HTTPS" {
		t.Fail()
	}
}
