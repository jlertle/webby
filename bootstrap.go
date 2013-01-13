package webby

type BootstrapHandler interface {
	Boot(*Web)
}

// Bootstrap Struct
type Bootstrap struct {
	functions []func(*Web)
}

// Register Functions to Bootstrap.
func (boot *Bootstrap) Register(functions ...func(*Web)) {
	boot.functions = append(boot.functions, functions...)
}

func NewBootstrapReg(functions ...func(*Web)) *Bootstrap {
	bo := &Bootstrap{}
	bo.Register(functions...)
	return bo
}

// Register Handler to Bootstrap.
func (boot *Bootstrap) RegisterHandler(handlers ...BootstrapHandler) {
	for _, handler := range handlers {
		ahandler := handler
		Function := func(w *Web) {
			ahandler.Boot(w)
		}
		boot.functions = append(boot.functions, Function)
	}
}

func NewBootstrapRegHandler(handlers ...BootstrapHandler) *Bootstrap {
	bo := &Bootstrap{}
	bo.RegisterHandler(handlers...)
	return bo
}

// Load Functions in Bootstrap.
func (boot *Bootstrap) Load(web *Web) {
	for _, function := range boot.functions {
		function(web)
		if web.CutOut() {
			return
		}
	}
}

var Boot = &Bootstrap{}

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
