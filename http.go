package webby

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
)

type Http struct {
	w *Web
}

func (w *Web) Http() Http {
	return Http{w}
}

// Set Cookie
func (h Http) SetCookie(cookie *http.Cookie) {
	http.SetCookie(h.w, cookie)
}

// Execute Handler
func (h Http) Exec(handler http.Handler) {
	handler.ServeHTTP(h.w, h.w.Req)
}

// Execute Function
func (h Http) ExecFunc(handler func(http.ResponseWriter, *http.Request)) {
	handler(h.w, h.w.Req)
}

// Get issues a GET to the specified URL.  If the response is one of the following
// redirect codes, Get follows the redirect, up to a maximum of 10 redirects:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// An error is returned if there were too many redirects or if there
// was an HTTP protocol error. A non-2xx response doesn't cause an
// error.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// Get is a wrapper around http.DefaultClient.Get.
func (h Http) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// Head issues a HEAD to the specified URL.  If the response is one of the
// following redirect codes, Head follows the redirect after calling the
// Client's CheckRedirect function.
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// Head is a wrapper around http.DefaultClient.Head
func (h Http) Head(url string) (*http.Response, error) {
	return http.Head(url)
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// Post is a wrapper around http.DefaultClient.Post
func (h Http) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, bodyType, body)
}

// PostForm issues a POST to the specified URL, with data's keys and
// values URL-encoded as the request body.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// PostForm is a wrapper around http.DefaultClient.PostForm
func (h Http) PostForm(url string, data url.Values) (*http.Response, error) {
	return http.PostForm(url, data)
}

// ReadResponse reads and returns an HTTP response from r.  The
// req parameter specifies the Request that corresponds to
// this Response.  Clients must call resp.Body.Close when finished
// reading resp.Body.  After that call, clients can inspect
// resp.Trailer to find key/value pairs included in the response
// trailer.
func (h Http) ReadResponse(r *bufio.Reader, req *http.Request) (*http.Response, error) {
	return http.ReadResponse(r, req)
}
