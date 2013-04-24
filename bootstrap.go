package webby

import (
	"sync"
)

type BootstrapHandler interface {
	Boot(*Web)
}

/*
Adapter for Converting Function to BootstrapHandler
*/
type FuncToBootstrapHandler func(*Web)

func (fn FuncToBootstrapHandler) Boot(w *Web) {
	fn(w)
}

// Bootstrap Struct, not to be confused with Twitter Bootstrap.
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
func (boot *Bootstrap) Register(functions ...FuncToBootstrapHandler) *Bootstrap {
	boot.Lock()
	defer boot.Unlock()
	for _, function := range functions {
		boot.boots = append(boot.boots, function)
	}
	return boot
}

// Construct New Bootstrap.
func NewBootstrap() *Bootstrap {
	return &Bootstrap{}
}

// Construct New Bootstrap and Register Functions.
func NewBootstrapReg(functions ...FuncToBootstrapHandler) *Bootstrap {
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
Bootstraps. (Not to be confused with Twitter Bootstrap)

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
