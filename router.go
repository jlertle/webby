package webby

import (
	"regexp"
	"sort"
	"sync"
)

type routerItem struct {
	RegExp         string
	RegExpComplied *regexp.Regexp
	Route          RouteHandler
}

// Route Map, use RegExp as key!
type RouteMap map[string]FuncToRouteHandler

// Route Handler Interface
type RouteHandler interface {
	View(*Web)
}

// Route Handler Map, use RegExp as key!
type RouteHandlerMap map[string]RouteHandler

type routes []*routerItem

func (ro routes) Len() int {
	return len(ro)
}

func (ro routes) Less(i, j int) bool {
	return ro[i].RegExp < ro[j].RegExp
}

func (ro routes) Swap(i, j int) {
	ro[i], ro[j] = ro[j], ro[i]
}

// Router (Controller), implement 'RouterHandler' interface
type Router struct {
	sync.RWMutex
	routes routes
}

func NewRouter() *Router {
	return &Router{}
}

func (ro *Router) getRoutes() routes {
	ro.RLock()
	defer ro.RUnlock()
	route := routes{}
	route = append(route, ro.routes...)
	return route
}

func (ro *Router) register(RegExpRule string, handler RouteHandler) {
	ro.Lock()
	defer ro.Unlock()
	for _, route := range ro.routes {
		if route.RegExp == RegExpRule {
			route.Route = handler
			return
		}
	}

	ro.routes = append(ro.routes, &routerItem{RegExpRule, regexp.MustCompile(RegExpRule), handler})
}

func (ro *Router) sortout() {
	ro.Lock()
	defer ro.Unlock()
	sort.Sort(ro.routes)
}

// Register rule and function to Router
func (ro *Router) Register(RegExpRule string, Function FuncToRouteHandler) *Router {
	ro.register(RegExpRule, Function)
	sort.Sort(ro.routes)
	return ro
}

// Register Map to Router
func (ro *Router) RegisterMap(routeMap RouteMap) *Router {
	if routeMap == nil {
		return ro
	}

	for rule, function := range routeMap {
		ro.register(rule, function)
	}
	ro.sortout()
	return ro
}

// Register rule and handler to Router
func (ro *Router) RegisterHandler(RegExpRule string, handler RouteHandler) *Router {
	ro.register(RegExpRule, handler)
	ro.sortout()
	return ro
}

// Register Handler Map to Router
func (ro *Router) RegisterHandlerMap(routeHandlerMap RouteHandlerMap) *Router {
	if routeHandlerMap == nil {
		return ro
	}

	for rule, handler := range routeHandlerMap {
		ro.register(rule, handler)
	}
	ro.sortout()
	return ro
}

func (ro *Router) load(w *Web, reset bool) bool {
	if reset {
		w.pri.path = w.Req.URL.Path
		w.pri.curpath = ""
	}

	for _, route := range ro.getRoutes() {
		if !route.RegExpComplied.MatchString(w.pri.path) {
			continue
		}

		w.pathDealer(route.RegExpComplied, pathStr(w.pri.path))

		w.RouteDealer(route.Route)
		return true
	}
	return false
}

func (ro *Router) debug(w *Web) {
	w.Status = 404
	out := w.Fmt()
	out.Print("404 Not Found\r\n\r\n")
	out.Print(w.Req.Host+w.pri.curpath, "\r\n\r\n")
	out.Print("Rule(s):\r\n")
	for _, route := range ro.getRoutes() {
		out.Print(route.RegExp, "\r\n")
	}
}

// Try to load matching route, output 404 on fail!
func (ro *Router) Load(w *Web) {
	if ro.load(w, false) {
		return
	}

	if w.Is().WebSocketRequest() {
		return
	}

	if DEBUG {
		ro.debug(w)
		return
	}

	w.Error404()
}

// Reset to root and try to load matching route, output 404 on fail!
func (ro *Router) LoadReset(w *Web) {
	if ro.load(w, true) {
		return
	}

	if w.Is().WebSocketRequest() {
		return
	}

	if DEBUG {
		ro.debug(w)
		return
	}

	w.Error404()
}

var Route = NewRouter()

// Router View
func (ro *Router) View(w *Web) {
	ro.Load(w)
}

// Implement RouteHandler interface!
type FuncToRouteHandler func(*Web)

func (fn FuncToRouteHandler) View(w *Web) {
	fn(w)
}

// Convert Router Handler to Function
func RouteHandlerToFunc(ro RouteHandler) func(w *Web) {
	aro := ro
	return func(w *Web) {
		aro.View(w)
	}
}

// Reset Url, Implement RouteHandler interface!
type RouteReset struct{ *Router }

func (ro RouteReset) View(w *Web) {
	ro.LoadReset(w)
}

func (w *Web) RouteDealer(ro RouteHandler) {
	for _, routeAssert := range _routeAsserter {
		if routeAssert.Assert(w, ro) {
			return
		}
	}

	switch t := ro.(type) {
	case MethodInterface:
		execMethodInterface(w, t)
	case ProtocolInterface:
		execProtocolInterface(w, t)
	default:
		ro.View(w)
	}
}

type RouteAsserter interface {
	Assert(*Web, RouteHandler) bool
}

type RouteAsserterFunc func(*Web, RouteHandler) bool

func (ra RouteAsserterFunc) Assert(w *Web, ro RouteHandler) bool {
	return ra(w, ro)
}

var _routeAsserter = []RouteAsserter{}

func RegisterRouteAsserter(ra ...RouteAsserter) {
	_routeAsserter = append(_routeAsserter, ra...)
}

func RegisterRouteAsserterFunc(ra ...RouteAsserterFunc) {
	for _, raa := range ra {
		RegisterRouteAsserter(raa)
	}
}
