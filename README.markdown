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

Webby here will only cover two components of the MVC, the View and the Controller, it's come with no Model, I decided not to include ORM because SQL is no longer the exclusive Database System as it once was years ago.  Basically use an existing dedicated ORM framework such as <https://github.com/astaxie/beedb> or <https://github.com/coopernurse/gorp>!  It should be easy to integrate into the web framework!