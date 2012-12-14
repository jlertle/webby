# Webby

Simple Web Application Server Framework!

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

	type index struct{}

	func (_ index) Socket(ws *websocket.Conn) {
		io.Copy(ws, ws)
	}

	func (_ index) View(w *webby.Web) {
		if w.IsWebSocketRequest() {
			websocket.Handler(func(ws *websocket.Conn) {
				index{}.Socket(ws)
			}).ServeHTTP(w.Res, w.Req)
			return
		}

		w.Print("<h1>Hello World!</h1>\r\n")
	}

	func init() {
		webby.Route.RegisterHandlerMap(webby.RouteHandlerMap{
			"^/$": index{},
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

Webby here will only cover two components of the MVC, the View and the Controller, it's come with no Model, I decided not to include ORM because SQL is no longer the exclusive Database System as it once was years ago.  Basically use an existing dedicated ORM framework such as <https://github.com/astaxie/beedb> or <https://github.com/coopernurse/gorp>!  It should be easy to integrate into the web framework!