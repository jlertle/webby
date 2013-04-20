/*
A Very Nice Micro Web Framework that can be used as a standalone or a embeddable solution for existing Go project that require a Web Interface!

It's built on top of the standard package 'net/http'!

Example:

	package main

	import (
		"code.google.com/p/go.net/websocket"
		"github.com/CJ-Jackson/webby"
		"io"
	)

	type Index struct {
		webby.Method
	}

	func (ind *Index) Get() {
		w := ind.W
		page := w.Param.GetInt("page")

		if page <= 0 {
			page = 1
		}

		w.Fmt().Print("<h1>Hello World!</h1>\r\n")
		w.Fmt().Print("Page: ", page, "\r\n")
	}

	func (ind *Index) Put() {
		data := struct {
			Title string `json:"title"`
			Life  int    `json:"life"`
		}{}

		// Decode from Request Body
		ind.W.Json().DecodeReqBody(&data)

		// Send it back to the client
		ind.W.Json().Send(data)

		// Rather use XML? Just replace 'Json' with 'Xml'!  Nice I know!
	}

	func (ind *Index) Delete() {
		// Delete a page, or anything really!
	}

	func (ind *Index) Ws() {
		if ind.W.Param.GetInt("page") > 0 {
			return
		}

		ind.W.RouteDealer(webby.HttpRouteHandler{websocket.Handler(
			func(ws *websocket.Conn) {
				io.Copy(ws, ws)
			})},
		)
	}

	func init() {
		webby.Route.RegisterHandlerMap(webby.RouteHandlerMap{
			// Main Route
			"^/$": &Index{},

			// Index Page Route
			"^/(?P<page>\\d+)/?$": &Index{},
		})
	}

	func main() {
		webby.Check(webby.StartHttp(":8080"))
	}

*/
package webby
