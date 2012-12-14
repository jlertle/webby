package webby

import (
	"regexp"
	"sort"
)

type routerItem struct {
	RegExp         string
	RegExpComplied *regexp.Regexp
	Function       func(*Web)
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

// Router (Controller)
type Router struct {
	routes routes
}

func (ro *Router) register(RegExpRule string, Function func(*Web)) {
	for _, route := range ro.routes {
		if route.RegExp == RegExpRule {
			route.Function = Function
			return
		}
	}

	ro.routes = append(ro.routes, &routerItem{RegExpRule, regexp.MustCompile(RegExpRule), Function})
}

// Register rule and function to Router
func (ro *Router) Register(RegExpRule string, Function func(*Web)) {
	ro.register(RegExpRule, Function)
	sort.Sort(ro.routes)
}

// Register Map to Router
func (ro *Router) RegisterMap(routeMap RouteMap) {
	if routeMap == nil {
		return
	}

	for rule, function := range routeMap {
		ro.register(rule, function)
	}
	sort.Sort(ro.routes)
}

func (ro *Router) registerHandler(RegExpRule string, handler RouteHandler) {
	ro.register(RegExpRule, func(w *Web) {
		handler.View(w)
	})
}

// Register rule and handler to Router
func (ro *Router) RegisterHandler(RegExpRule string, handler RouteHandler) {
	ro.registerHandler(RegExpRule, handler)
	sort.Sort(ro.routes)
}

// Register Handler Map to Router
func (ro *Router) RegisterHandlerMap(routeHandlerMap RouteHandlerMap) {
	if routeHandlerMap == nil {
		return
	}

	for rule, handler := range routeHandlerMap {
		ro.registerHandler(rule, handler)
	}
	sort.Sort(ro.routes)
}

func (ro *Router) load(w *Web, reset bool) bool {
	if reset {
		w.path = w.Req.URL.Path
		w.curpath = ""
	}

	for _, route := range ro.routes {
		if route.RegExpComplied.MatchString(w.path) {
			names := route.RegExpComplied.SubexpNames()
			matches := route.RegExpComplied.FindStringSubmatch(w.path)

			w.curpath += matches[0]

			for key, name := range names {
				if name != "" {
					str := matches[key]
					num, err := toInt(str)
					if err == nil {
						w.Param[name] = num
					} else {
						w.Param[name] = str
					}
				}
			}

			w.path = route.RegExpComplied.ReplaceAllString(w.path, "")

			route.Function(w)
			return true

		}
	}
	return false
}

// Try to load matching route, output 404 on fail!
func (ro *Router) Load(w *Web) {
	if ro.load(w, false) {
		return
	}

	if w.IsWebSocketRequest() {
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

	Error404(w)
}

var Route = &Router{}
