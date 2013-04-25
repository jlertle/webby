package webby

import (
	"reflect"
	"sync"
)

// Middleware Interface
type MiddlewareInterface interface {
	Init(*Web)
	Pre()
	Post()
}

// Implement MiddlewareInterface
type Middleware struct {
	W *Web
}

// Init
func (mid *Middleware) Init(w *Web) {
	mid.W = w
}

// Pre boot
func (mid *Middleware) Pre() {
	// Do nothing
}

// Post boot
func (mid *Middleware) Post() {
	// Do nothing
}

type Middlewares struct {
	sync.Mutex
	items []MiddlewareInterface
	w     *Web
}

// Construct New Middleware
func NewMiddlewares() *Middlewares {
	return &Middlewares{}
}

// Register Middlewares
func (mid *Middlewares) Register(middlewares ...MiddlewareInterface) *Middlewares {
	if mid.w == nil {
		mid.Lock()
		defer mid.Unlock()
	}
	if mid.items == nil {
		mid.items = []MiddlewareInterface{}
	}
	mid.items = append(mid.items, middlewares...)
	return mid
}

func (mid *Middlewares) getItems() []MiddlewareInterface {
	mid.Lock()
	defer mid.Unlock()
	items := []MiddlewareInterface{}
	items = append(items, mid.items...)
	return items
}

// Init Middlewares, return initialised structure. 
func (mid *Middlewares) Init(w *Web) *Middlewares {
	if mid.w != nil {
		return mid
	}
	middlewares := NewMiddlewares()
	middlewares.w = w
	for _, middleware := range mid.getItems() {
		newmiddleware := reflect.New(
			reflect.Indirect(reflect.ValueOf(middleware)).Type()).Interface().(MiddlewareInterface)
		newmiddleware.Init(w)
		middlewares.Register(newmiddleware)
	}
	return middlewares
}

// Pre boot
func (mid *Middlewares) Pre() {
	if mid.w == nil {
		return
	}
	for _, middleware := range mid.items {
		middleware.Pre()
		if mid.w.CutOut() {
			return
		}
	}
}

// Post boot, you may want to use the keyword 'defer'
func (mid *Middlewares) Post() {
	if mid.w == nil {
		return
	}
	for _, middleware := range mid.items {
		middleware.Post()
	}
}

// Default Middleware Holders!
var (
	MainMiddlewares = NewMiddlewares()
	AppMiddlewares  = NewMiddlewares()
)
