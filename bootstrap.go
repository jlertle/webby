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

func NewBootstrap() *Bootstrap {
	return &Bootstrap{}
}

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

// For Allowing Libraries to Automatically Plugin into this Framework. 
var MainBoot = NewBootstrap()

// For Allowing Web Application to Add Function the Framework.
var Boot = NewBootstrap()

// For Allowing Libraries to Add Function to the Html Template Engine!
var HtmlFuncBoot = NewBootstrap()

// Implement RouteHandler interface.
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

type BootRoutePipe struct {
	br BootRoute
}

func NewBootRoute() BootRoutePipe {
	return BootRoutePipe{
		br: BootRoute{},
	}
}

func (bo BootRoutePipe) Boot(boot *Bootstrap) BootRoutePipe {
	bo.br.BOOT = boot
	return bo
}

func (bo BootRoutePipe) Router(router *Router) BootRoutePipe {
	bo.br.ROUTER = router
	return bo
}

func (bo BootRoutePipe) Get() BootRoute {
	return bo.br
}
