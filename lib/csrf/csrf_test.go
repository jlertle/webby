package csrf

import (
	"github.com/CJ-Jackson/webby"
)

func ExampleCheck() {
	// Registering Bootsrap
	webby.Boot.RegisterHandler(Check{})
}
