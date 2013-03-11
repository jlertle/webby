package webby

import (
	"time"
)

func CurTime() time.Time {
	return time.Now()
}

func (w *Web) CurTime() time.Time {
	return CurTime()
}
