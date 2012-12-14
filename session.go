package webby

import (
	"fmt"
	"net/http"
	"time"
)

var (
	SessionCookieName          = "__session"
	SessionEnabled             = true
	SessionExpire              = 20 * time.Minute
	SessionExpiryCheckInterval = 10 * time.Minute
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

var session_map = map[string]sessionInterface{}

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

// Set Session
func (w *Web) SetSession(data interface{}) {
	if !SessionEnabled {
		return
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
	session_map[sesCookie.Value] = &session{data, time.Now().Add(SessionExpire)}
}

// Init Session
func (w *Web) initSession() {
	if !SessionEnabled {
		return
	}

	sesCookie, err := w.GetCookie(SessionCookieName)
	if err != nil {
		return
	}

	switch t := session_map[sesCookie.Value].(type) {
	case *session:
		if time.Now().Unix() < t.getExpire().Unix() {
			w.Session = t.getData()
			t.hit()
			return
		}
	}

	sesCookie.MaxAge = -1
	w.SetCookie(sesCookie)
}

// Destroy Session
func (w *Web) DestroySession() {
	if !SessionEnabled {
		return
	}

	sesCookie, err := w.GetCookie(SessionCookieName)
	if err != nil {
		return
	}

	switch session_map[sesCookie.Value].(type) {
	case *session:
		delete(session_map, sesCookie.Value)
	}
	sesCookie.MaxAge = -1
	w.SetCookie(sesCookie)
}

//	Session Expiry Check in a loop
func sessionExpiryCheck() {
	for {
		time.Sleep(SessionExpiryCheckInterval)
		if !SessionEnabled {
			break
		}
		curtime := time.Now()
		for key, value := range session_map {
			if curtime.Unix() > value.getExpire().Unix() {
				delete(session_map, key)
			}
		}
	}
}

//	Start Session Check
func init() {
	go sessionExpiryCheck()
}
