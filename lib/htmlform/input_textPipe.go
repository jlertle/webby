package htmlform

type PipeInputText struct {
	in InputText
}

func NewInputText(name string) PipeInputText {
	return PipeInputText{
		InputText{Name: name},
	}
}

func (in PipeInputText) Get() *InputText {
	return &in.in
}

func (in PipeInputText) Value(value string) PipeInputText {
	in.in.Value = value
	return in
}

func (in PipeInputText) Id(id string) PipeInputText {
	in.in.Id = id
	return in
}

func (in PipeInputText) Class(class string) PipeInputText {
	in.in.Class = class
	return in
}

func (in PipeInputText) MinChar(minchar int) PipeInputText {
	in.in.MinChar = minchar
	return in
}

func (in PipeInputText) Mandatory() PipeInputText {
	return in.MinChar(1)
}

func (in PipeInputText) MaxChar(maxchar int) PipeInputText {
	in.in.MaxChar = maxchar
	return in
}

func (in PipeInputText) RegExp(rule, errMsg string) PipeInputText {
	in.in.RegExpRule = rule
	in.in.RegExpErr = errMsg
	return in
}

func (in PipeInputText) MustMatch(match, errMsg string) PipeInputText {
	in.in.MustMatch = match
	in.in.MustMatchErr = errMsg
	return in
}

func (in PipeInputText) Extra(extra func(*Validation) error) PipeInputText {
	in.in.extra = extra
	return in
}
