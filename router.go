package webby

import (
	"regexp"
	"sort"
	"sync"
)

type Param map[string]string

func (pa Param) Add(name, value string) {
	pa[name] = value
}

func (pa Param) Get(name string) string {
	return pa[name]
}

func (pa Param) GetInt64(name string) int64 {
	num := int64(0)
	var err error
	num, err = toInt(pa[name])
	if err != nil {
		return 0
	}
	return num
}

func (pa Param) GetInt(name string) int {
	return int(pa.GetInt64(name))
}

func (pa Param) GetUint64(name string) uint64 {
	num := uint64(0)
	var err error
	num, err = toUint(pa[name])
	if err != nil {
		return 0
	}
	return num
}

func (pa Param) GetUint(name string) uint {
	return uint(pa.GetUint64(name))
}

func (pa Param) GetFloat64(name string) float64 {
	num := float64(0)
	var err error
	num, err = toFloat(pa[name])
	if err != nil {
		return float64(0)
	}
	return num
}

func (pa Param) GetFloat32(name string) float32 {
	return float32(pa.GetFloat64(name))
}

type routerItem struct {
	RegExp         string
	RegExpComplied *regexp.Regexp
	Route          RouteHandler
}

// Route Map, use RegExp as key!
type RouteMap map[string]func(*Web)

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
func (ro *Router) Register(RegExpRule string, Function func(*Web)) {
	ro.register(RegExpRule, FuncToRouteHandler{Function})
	sort.Sort(ro.routes)
}

// Register Map to Router
func (ro *Router) RegisterMap(routeMap RouteMap) {
	if routeMap == nil {
		return
	}

	for rule, function := range routeMap {
		ro.register(rule, FuncToRouteHandler{function})
	}
	ro.sortout()
}

func NewRouterMap(routeMap RouteMap) *Router {
	ro := &Router{}
	ro.RegisterMap(routeMap)
	return ro
}

// Register rule and handler to Router
func (ro *Router) RegisterHandler(RegExpRule string, handler RouteHandler) {
	ro.register(RegExpRule, handler)
	ro.sortout()
}

// Register Handler Map to Router
func (ro *Router) RegisterHandlerMap(routeHandlerMap RouteHandlerMap) {
	if routeHandlerMap == nil {
		return
	}

	for rule, handler := range routeHandlerMap {
		ro.register(rule, handler)
	}
	ro.sortout()
}

func NewRouterHandlerMap(routeHandlerMap RouteHandlerMap) *Router {
	ro := &Router{}
	ro.RegisterHandlerMap(routeHandlerMap)
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

		names := route.RegExpComplied.SubexpNames()
		matches := route.RegExpComplied.FindStringSubmatch(w.pri.path)

		w.pri.curpath += matches[0]

		for key, name := range names {
			if name != "" {
				w.Param.Add(name, matches[key])
			}
		}

		w.pri.path = w.pri.path[route.RegExpComplied.FindStringIndex(w.pri.path)[1]:]

		route.Route.View(w)
		return true
	}
	return false
}

func (ro *Router) debug(w *Web) {
	w.Status = 404
	w.Print("404 Not Found\r\n\r\n")
	w.Print(w.Req.Host+w.pri.curpath, "\r\n\r\n")
	w.Print("Rule(s):\r\n")
	for _, route := range ro.getRoutes() {
		w.Print(route.RegExp, "\r\n")
	}
}

// Try to load matching route, output 404 on fail!
func (ro *Router) Load(w *Web) {
	if ro.load(w, false) {
		return
	}

	if w.IsWebSocketRequest() {
		return
	}

	if DEBUG {
		ro.debug(w)
		return
	}

	Error404(w)
}

// Reset to root and try to load matching route, output 404 on fail!
func (ro *Router) LoadReset(w *Web) {
	if ro.load(w, true) {
		return
	}

	if w.IsWebSocketRequest() {
		return
	}

	if DEBUG {
		ro.debug(w)
		return
	}

	Error404(w)
}

var Route = &Router{}

// Router View
func (ro *Router) View(w *Web) {
	ro.Load(w)
}

// Implement RouteHandler interface!
type FuncToRouteHandler struct {
	Function func(*Web)
}

func (fn FuncToRouteHandler) View(w *Web) {
	fn.Function(w)
}

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
