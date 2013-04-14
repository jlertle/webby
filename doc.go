/*
A Very Nice Micro Web Framework that can be used as a standalone or a embeddable solution for existing Go project that require a Web Interface!

It's built on top of the standard package 'net/http'!

Example:

	package main

	import (
		"code.google.com/p/go.net/websocket"
		"github.com/cj-jackson/webby"
		"io"
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
		webby.Route.RegisterHandlerMap(webby.RouteHandlerMap{
			// Main Route
			"^/$": webby.NewJunction().Websocket(
				webby.HttpRouteHandler{websocket.Handler(EchoServer)},
			).Get(index{}).GetJunction(),

			// Index Page Route
			"^/(?P<page>\\d+)/?$": index{},
		})
	}

	func main() {
		webby.Check(webby.StartHttp(":8080"))
	}

*/
package webby
