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
		boot.functions = append(boot.functions, func(w *Web) {
			handler.Boot(w)
		})
	}
}

// Load Functions in Bootstrap.
func (boot *Bootstrap) Load(web *Web) {
	for _, function := range boot.functions {
		function(web)
	}
}

var Boot = &Bootstrap{}
