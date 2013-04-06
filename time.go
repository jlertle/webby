package webby

import (
	"time"
)

// Default Time Location
var DefaultTimeLoc *time.Location

// Default Time Format
var DefaultTimeFormat string = "Monday, _2 January 2006, 15:04"

// Set Timezone on global level
func SetTimeZone(zone string) {
	var err error
	DefaultTimeLoc, err = time.LoadLocation(zone)
	Check(err)
}

// Set Timezone on user request level
func (w *Web) SetTimeZone(zone string) {
	var err error
	w.TimeLoc, err = time.LoadLocation(zone)
	w.Check(err)
}

// Get Current Time
func CurTime() time.Time {
	return time.Now()
}

// Get Current Time
func (w *Web) CurTime() time.Time {
	return CurTime()
}

func init() {
	DefaultTimeLoc, _ = time.LoadLocation("Local")

	HtmlFuncBoot.Register(func(w *Web) {
		// Convert to Default Timezone.
		w.HtmlFunc["time"] = func(clock time.Time) time.Time {
			return clock.In(w.TimeLoc)
		}

		// Convert to Timezone
		w.HtmlFunc["timeZone"] = func(zone string, clock time.Time) time.Time {
			loc, err := time.LoadLocation(zone)
			w.Check(err)
			return clock.In(loc)
		}

		// Format time, leave empty for default
		w.HtmlFunc["timeFormat"] = func(format string, clock time.Time) string {
			if format == "" {
				format = w.TimeFormat
			}
			return clock.Format(format)
		}
	})
}
