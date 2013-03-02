// HTML Form Helper Internationalization Library!
package lang

import (
	"encoding/gob"
	"sync"
)

type Lang map[string]string

type LangLocker struct {
	sync.RWMutex
	langs map[string]Lang
}

func (lg *LangLocker) Add(name string, lang Lang) {
	if lg.langs == nil {
		lg.langs = map[string]Lang{}
	}
	lg.langs[name] = lang
}

func (lg *LangLocker) Get(name string) Lang {
	lg.RLock()
	defer lg.RUnlock()
	if lg.langs == nil {
		return nil
	}
	if lg.langs[name] == nil {
		return nil
	}
	lang := Lang{}
	for name, value := range lg.langs[name] {
		lang[name] = value
	}
	return lang
}

var Langs = &LangLocker{}

func init() {
	gob.Register(Lang{})
}
