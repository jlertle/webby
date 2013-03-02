package webby

import (
	"sync"
)

type BootstrapHandler interface {
	Boot(*Web)
}

type FuncToBootstrapHandler struct {
	function func(*Web)
}

func (fn FuncToBootstrapHandler) Boot(w *Web) {
	fn.function(w)
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
func (boot *Bootstrap) Register(functions ...func(*Web)) {
	boot.Lock()
	defer boot.Unlock()
	for _, function := range functions {
		boot.boots = append(boot.boots, FuncToBootstrapHandler{function})
	}
}

func NewBootstrapReg(functions ...func(*Web)) *Bootstrap {
	bo := &Bootstrap{}
	bo.Register(functions...)
	return bo
}

// Register Handler to Bootstrap.
func (boot *Bootstrap) RegisterHandler(handlers ...BootstrapHandler) {
	boot.Lock()
	defer boot.Unlock()
	boot.boots = append(boot.boots, handlers...)
}

func NewBootstrapRegHandler(handlers ...BootstrapHandler) *Bootstrap {
	bo := &Bootstrap{}
	bo.RegisterHandler(handlers...)
	return bo
}

// Load Functions in Bootstrap.
func (boot *Bootstrap) Load(web *Web) {
	for _, bt := range boot.getBoots() {
		bt.Boot(web)
		if web.CutOut() {
			return
		}
	}
}

// For Allowing Libraries to Automatically Plugin into this Framework. 
var MainBoot = &Bootstrap{}

// For Allowing Web Application to Add Function the Framework.
var Boot = &Bootstrap{}

// For Allowing Libraries to Add Function to the Html Template Engine!
var HtmlFuncBoot = &Bootstrap{}

type BootRoute struct {
	BOOT   *Bootstrap
	ROUTER *Router
}

func (bo BootRoute) View(web *Web) {
	if bo.BOOT != nil {
		bo.BOOT.Load(web)

		if web.CutOut() {
			return
		}
	}

	if bo.ROUTER != nil {
		bo.ROUTER.LoadReset(web)

		if web.CutOut() {
			return
		}
	}
}
