package webby

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	SessionCookieName          = "__session"
	SessionExpire              = 20 * time.Minute
	SessionExpiryCheckInterval = 10 * time.Minute
	sessionExpiryCheckActive   = false
)

// Structure of Session
type session struct {
	Data   interface{}
	Expire time.Time
}

// Get the session data!
func (ses *session) getData() interface{} {
	return ses.Data
}

// Returns the time of expiry
func (ses *session) getExpire() time.Time {
	return ses.Expire
}

// Reset Expiry Time to 20 minutes in advanced!
func (ses *session) hit() {
	ses.Expire = time.Now().Add(SessionExpire)
}

// Session interface
//
// Note: Useful for checking the existant of the session!
type sessionInterface interface {
	getData() interface{}
	getExpire() time.Time
	hit()
}

var sessionMap = struct {
	sync.Mutex
	m map[string]sessionInterface
}{m: map[string]sessionInterface{}}

// Convert Unsigned 64-bit Int to Bytes.
func uint64ToByte(num uint64) [8]byte {
	var buf [8]byte
	buf[0] = byte(num >> 0)
	buf[1] = byte(num >> 8)
	buf[2] = byte(num >> 16)
	buf[3] = byte(num >> 24)
	buf[4] = byte(num >> 32)
	buf[5] = byte(num >> 40)
	buf[6] = byte(num >> 48)
	buf[7] = byte(num >> 56)
	return buf
}

type SessionHandler interface {
	Set(*Web, interface{})
	Init(*Web)
	Destroy(*Web)
}

type SessionMemory struct{}

func (_ SessionMemory) Set(w *Web, data interface{}) {
	sessionMap.Lock()
	defer sessionMap.Unlock()

	if !sessionExpiryCheckActive {
		sessionExpiryCheckActive = true
		go sessionExpiryCheck()
	}

	sesCookie, err := w.GetCookie(SessionCookieName)

	if err != nil {
		sesCookie = &http.Cookie{}
		curtime := time.Now()
		sesCookie.Name = SessionCookieName
		sesCookie.Value = fmt.Sprintf("%x%x", uint64ToByte(uint64(curtime.Unix())),
			uint64ToByte(uint64(curtime.UnixNano())))
	}

	w.SetCookie(sesCookie)
	sessionMap.m[sesCookie.Value] = &session{data, time.Now().Add(SessionExpire)}
}

func deleteSessionFromMap(key string) {
	delete(sessionMap.m, key)
}

func (_ SessionMemory) Init(w *Web) {
	sessionMap.Lock()
	defer sessionMap.Unlock()

	sesCookie, err := w.GetCookie(SessionCookieName)
	if err != nil {
		return
	}

	switch t := sessionMap.m[sesCookie.Value].(type) {
	case *session:
		if time.Now().Unix() < t.getExpire().Unix() {
			w.Session = t.getData()
			t.hit()
			return
		}
	}

	deleteSessionFromMap(sesCookie.Value)
	sesCookie.MaxAge = -1
	w.SetCookie(sesCookie)
}

func (_ SessionMemory) Destroy(w *Web) {
	sessionMap.Lock()
	defer sessionMap.Unlock()

	sesCookie, err := w.GetCookie(SessionCookieName)
	if err != nil {
		return
	}

	switch sessionMap.m[sesCookie.Value].(type) {
	case *session:
		deleteSessionFromMap(sesCookie.Value)
	}
	sesCookie.MaxAge = -1
	w.SetCookie(sesCookie)
}

const sessionFileExt = ".wbs"

type SessionFile struct {
	Path string
}

func (se SessionFile) Set(w *Web, data interface{}) {
	sesCookie, err := w.GetCookie(SessionCookieName)

	if err != nil {
		sesCookie = &http.Cookie{}
		curtime := time.Now()
		sesCookie.Name = SessionCookieName
		sesCookie.Value = fmt.Sprintf("%x%x", uint64ToByte(uint64(curtime.Unix())),
			uint64ToByte(uint64(curtime.UnixNano())))
	}

	w.SetCookie(sesCookie)
	file, err := os.Create(se.Path + "/" + sesCookie.Value + sessionFileExt)
	w.Check(err)

	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(&session{data, time.Now().Add(SessionExpire)})
	if err != nil {
		panic(err)
	}
}

func (se SessionFile) Init(w *Web) {
	sesCookie, err := w.GetCookie(SessionCookieName)
	if err != nil {
		return
	}

	file, err := os.Open(se.Path + "/" + sesCookie.Value + sessionFileExt)
	if err != nil {
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)

	ses := &session{}

	err = dec.Decode(&ses)
	if err != nil {
		panic(err)
	}

	if time.Now().Unix() < ses.getExpire().Unix() {
		w.Session = ses.getData()
		ses.hit()
		return
	}

	os.Remove(se.Path + "/" + sesCookie.Value + sessionFileExt)
	sesCookie.MaxAge = -1
	w.SetCookie(sesCookie)
}

func (se SessionFile) Destroy(w *Web) {
	sesCookie, err := w.GetCookie(SessionCookieName)
	if err != nil {
		return
	}

	os.Remove(se.Path + "/" + sesCookie.Value + sessionFileExt)
	sesCookie.MaxAge = -1
	w.SetCookie(sesCookie)
}

var DefaultSessionHandler SessionHandler = SessionMemory{}

// Set Session
func (w *Web) SetSession(data interface{}) {
	DefaultSessionHandler.Set(w, data)
}

// Init Session
func (w *Web) initSession() {
	DefaultSessionHandler.Init(w)
}

// Destroy Session
func (w *Web) DestroySession() {
	DefaultSessionHandler.Destroy(w)
}

//	Session Expiry Check in a loop
func sessionExpiryCheck() {
	for {
		time.Sleep(SessionExpiryCheckInterval)
		curtime := time.Now()

		sessionMap.Lock()

		if len(sessionMap.m) <= 0 {
			sessionExpiryCheckActive = false
			sessionMap.Unlock()
			break
		}
		for key, value := range sessionMap.m {
			if curtime.Unix() > value.getExpire().Unix() {
				delete(sessionMap.m, key)
			}
		}

		sessionMap.Unlock()
	}
}
