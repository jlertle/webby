package webby

import (
	"encoding/gob"
)

type sessionAdvMap map[string]interface{}

func init() {
	gob.Register(sessionAdvMap{})
}

// Key/Value version of Session
type SessionAdv struct {
	w    *Web
	sMap sessionAdvMap
}

// Key/Value version of Session
func (w *Web) SessionAdv() *SessionAdv {
	if w.pri.session != nil {
		return w.pri.session
	}
	w.pri.session = &SessionAdv{w: w}
	w.pri.session.init()
	return w.pri.session
}

func (se *SessionAdv) init() {
	if se.sMap != nil {
		return
	}

	switch t := se.w.Session.(type) {
	case sessionAdvMap:
		se.sMap = t
	default:
		se.sMap = sessionAdvMap{}
	}
}

// Set Session by Key
func (se *SessionAdv) Set(key string, value interface{}) {
	se.sMap[key] = value
}

// Get Session by Key
func (se *SessionAdv) Get(key string) interface{} {
	return se.sMap[key]
}

// Save Session and set Cookie to client.
func (se *SessionAdv) Save() {
	se.w.SetSession(se.sMap)
}
