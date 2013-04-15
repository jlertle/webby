package webby

type Cmd struct {
	w *Web
}

func (w *Web) Cmd() Cmd {
	return Cmd{w}
}

// Set Custom Command
func (c Cmd) SetCmd(name string, cmd func(interface{}) interface{}) {
	c.w.pri.cmd[name] = cmd
}

// Execute Custom Command
func (c Cmd) ExecCmd(name string, v interface{}) interface{} {
	if c.w.pri.cmd[name] == nil {
		panic(ErrorStr("CMD: '" + name + "' does not exist!"))
	}
	return c.w.pri.cmd[name](v)
}
