package htmlform

type PipeInputRadio struct {
	ra InputRadio
}

func NewInputRadio(name string) PipeInputRadio {
	return PipeInputRadio{InputRadio{Name: name}}
}

func (ra PipeInputRadio) Get() *InputRadio {
	return &ra.ra
}

func (ra PipeInputRadio) Value(value string) PipeInputRadio {
	ra.ra.Value = value
	return ra
}

func (ra PipeInputRadio) Id(id string) PipeInputRadio {
	ra.ra.Id = id
	return ra
}

func (ra PipeInputRadio) Class(class string) PipeInputRadio {
	ra.ra.Class = class
	return ra
}

func (ra PipeInputRadio) Selected() PipeInputRadio {
	ra.ra.Selected = true
	return ra
}
