package webby

import (
	"time"
)

func CurTime() time.Time {
	return time.Now()
}

func (w *Web) Curtime() time.Time {
	return CurTime()
}
