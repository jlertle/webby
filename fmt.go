package webby

import (
	"fmt"
)

type Fmt struct {
	w *Web
}

func (w *Web) Fmt() Fmt {
	return Fmt{w}
}

// Print formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(f.w, a...)
}

// Printf formats according to a format specifier and writes to client (http web server or browser).
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(f.w, format, a...)
}

// Println formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(f.w, a...)
}
