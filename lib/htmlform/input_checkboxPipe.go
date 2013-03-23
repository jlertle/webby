package htmlform

type PipeInputCheckbox struct {
	ch InputCheckbox
}

func NewInputCheckbox(name string) PipeInputCheckbox {
	return PipeInputCheckbox{InputCheckbox{Name: name}}
}

func (ch PipeInputCheckbox) Get() *InputCheckbox {
	return &ch.ch
}

func (ch PipeInputCheckbox) Value(value string) PipeInputCheckbox {
	ch.ch.Value = value
	return ch
}

func (ch PipeInputCheckbox) Id(id string) PipeInputCheckbox {
	ch.ch.Id = id
	return ch
}

func (ch PipeInputCheckbox) Class(class string) PipeInputCheckbox {
	ch.ch.Class = class
	return ch
}

func (ch PipeInputCheckbox) Selected() PipeInputCheckbox {
	ch.ch.Selected = true
	return ch
}

func (ch PipeInputCheckbox) Mandatory() PipeInputCheckbox {
	ch.ch.Mandatory = true
	return ch
}
