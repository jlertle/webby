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
