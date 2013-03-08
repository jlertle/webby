# Webby

A Very Nice Web Framework that can be used as a standalone or a embeddable solution for existing Go project that require a Web Interface!

It's built on top of the standard package 'net/http'!

## Installation ##

	go get github.com/CJ-Jackson/webby
	
## Example Web Application ##

	package main

	import (
		"code.google.com/p/go.net/websocket"
		"github.com/CJ-Jackson/webby"
		"io"
		"log"
		"net/http"
	)

	func EchoServer(ws *websocket.Conn) {
		io.Copy(ws, ws)
	}

	type index struct{}

	func (_ index) View(w *webby.Web) {
		page := w.Param.GetInt("page")

		if page <= 0 {
			page = 1
		}

		w.Print("<h1>Hello World!</h1>\r\n")
		w.Print("Page: ", page, "\r\n")
	}

	func init() {
		index_junc := webby.Junction{
			WS:  webby.HttpRouteHandler{websocket.Handler(EchoServer)},
			ALL: index{},
		}

		webby.Route.RegisterHandlerMap(webby.RouteHandlerMap{
			"^/$": index_junc,
			"^/(?P<page>\\d+)/?$": index_junc,
		})
	}

	func main() {
		err := http.ListenAndServe(":8080", webby.Web{})

		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}


Note: This example require websocket library installed; The framework itself can live without that.

	go get code.google.com/p/go.net/websocket

## Documentation ##

Documentation are avaliable at

http://go.pkgdoc.org/github.com/CJ-Jackson/webby

## Note ##

Note this framework does not have an ORM but can co-exist with an ORM Framework, a Google App Engine user probably won't need an ORM Framework!