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
	post  bool
}

// Construct New Bootstrap.
func NewBootstrap() *Bootstrap {
	return &Bootstrap{post: false}
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

// Register Handler to Bootstrap.
func (boot *Bootstrap) RegisterHandler(handlers ...BootstrapHandler) *Bootstrap {
	boot.Lock()
	defer boot.Unlock()
	boot.boots = append(boot.boots, handlers...)
	return boot
}

// Load Functions in Bootstrap.
func (boot *Bootstrap) Load(w *Web) {
	for _, bt := range boot.getBoots() {
		bt.Boot(w)
		if w.CutOut() && !boot.post {
			return
		}
	}
}

func (boot *Bootstrap) PostMode() *Bootstrap {
	boot.post = true
	return boot
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

You can think of them as filters.
*/
var (
	MainBoot     = NewBootstrap()
	Boot         = NewBootstrap()
	HtmlFuncBoot = NewBootstrap()
	PostBoot     = NewBootstrap().PostMode()
)
