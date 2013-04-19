package webby

import (
	"sync"
)

type BootstrapHandler interface {
	Boot(*Web)
}

type FuncToBootstrapHandler struct {
	Function func(*Web)
}

func (fn FuncToBootstrapHandler) Boot(w *Web) {
	fn.Function(w)
}

// Bootstrap Struct
type Bootstrap struct {
	sync.RWMutex
	boots []BootstrapHandler
}

func (boot *Bootstrap) getBoots() []BootstrapHandler {
	boot.RLock()
	defer boot.RUnlock()
	boots := []BootstrapHandler{}
	boots = append(boots, boot.boots...)
	return boots
}

// Register Functions to Bootstrap.
func (boot *Bootstrap) Register(functions ...func(*Web)) *Bootstrap {
	boot.Lock()
	defer boot.Unlock()
	for _, function := range functions {
		boot.boots = append(boot.boots, FuncToBootstrapHandler{function})
	}
	return boot
}

// Construct New Bootstrap.
func NewBootstrap() *Bootstrap {
	return &Bootstrap{}
}

// Construct New Bootstrap and Register Functions.
func NewBootstrapReg(functions ...func(*Web)) *Bootstrap {
	bo := &Bootstrap{}
	bo.Register(functions...)
	return bo
}

// Register Handler to Bootstrap.
func (boot *Bootstrap) RegisterHandler(handlers ...BootstrapHandler) *Bootstrap {
	boot.Lock()
	defer boot.Unlock()
	boot.boots = append(boot.boots, handlers...)
	return boot
}

// Construct New Bootstrap and Register Handlers.
func NewBootstrapRegHandler(handlers ...BootstrapHandler) *Bootstrap {
	bo := &Bootstrap{}
	bo.RegisterHandler(handlers...)
	return bo
}

// Load Functions in Bootstrap.
func (boot *Bootstrap) Load(w *Web) {
	for _, bt := range boot.getBoots() {
		bt.Boot(w)
		if w.CutOut() {
			return
		}
	}
}

func (boot *Bootstrap) Boot(w *Web) {
	boot.Load(w)
}

/*
Bootstraps.

MainBoot is framework level, you can use it for something like stripping www from the url.

Boot is application level.

HtmlFuncBoot is exclusively used for adding functions to html template engine.

PostBoot is the last to be executed, do not write output to client at that point, it's for things like logging!

You can think of them as middle-ware or filters.
*/
var (
	MainBoot     = NewBootstrap()
	Boot         = NewBootstrap()
	HtmlFuncBoot = NewBootstrap()
	PostBoot     = NewBootstrap()
)

// Bootstrap and Router wrapper. Implement RouteHandler interface.
type BootRoute struct {
	BOOT   *Bootstrap
	ROUTER *Router
}

func (bo BootRoute) View(w *Web) {
	if bo.BOOT != nil {
		bo.BOOT.Load(w)

		if w.CutOut() {
			return
		}
	}

	if bo.ROUTER != nil {
		bo.ROUTER.LoadReset(w)

		if w.CutOut() {
			return
		}
	}
}

// Chainable version of BootRoute
type PipeBootRoute struct {
	br BootRoute
}

// PipeBootRoute constructor
func NewBootRoute() PipeBootRoute {
	return PipeBootRoute{
		br: BootRoute{},
	}
}

// Set boot
func (bo PipeBootRoute) Boot(boot *Bootstrap) PipeBootRoute {
	bo.br.BOOT = boot
	return bo
}

// Set Router
func (bo PipeBootRoute) Router(router *Router) PipeBootRoute {
	bo.br.ROUTER = router
	return bo
}

// Get BootRoute
func (bo PipeBootRoute) Get() BootRoute {
	return bo.br
}
