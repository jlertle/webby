package htmlform

type PipeInputPassword struct {
	in InputPassword
}

func NewInputPassword(name string) PipeInputPassword {
	return PipeInputPassword{
		InputPassword{Name: name},
	}
}

func (in PipeInputPassword) Get() *InputPassword {
	return &in.in
}

func (in PipeInputPassword) GetStruct() FormHandlerExt {
	return in.Get()
}

func (in PipeInputPassword) Value(value string) PipeInputPassword {
	in.in.Value = value
	return in
}

func (in PipeInputPassword) Id(id string) PipeInputPassword {
	in.in.Id = id
	return in
}

func (in PipeInputPassword) Class(class string) PipeInputPassword {
	in.in.Class = class
	return in
}

func (in PipeInputPassword) MinChar(minchar int) PipeInputPassword {
	in.in.MinChar = minchar
	return in
}

func (in PipeInputPassword) Mandatory() PipeInputPassword {
	return in.MinChar(1)
}

func (in PipeInputPassword) MaxChar(maxchar int) PipeInputPassword {
	in.in.MaxChar = maxchar
	return in
}

func (in PipeInputPassword) RegExp(rule, errMsg string) PipeInputPassword {
	in.in.RegExpRule = rule
	in.in.RegExpErr = errMsg
	return in
}

func (in PipeInputPassword) MustMatch(match, errMsg string) PipeInputPassword {
	in.in.MustMatch = match
	in.in.MustMatchErr = errMsg
	return in
}

func (in PipeInputPassword) Placeholder(placeholder string) PipeInputPassword {
	in.in.Placeholder = placeholder
	return in
}

func (in PipeInputPassword) Extra(extra ExtraFunc) PipeInputPassword {
	in.in.extra = extra
	return in
}
