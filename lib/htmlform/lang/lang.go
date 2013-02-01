// HTML Form Helper Language Library!
package lang

import (
	"encoding/gob"
)

type Lang map[string]string

func init() {
	gob.Register(Lang{})
}
