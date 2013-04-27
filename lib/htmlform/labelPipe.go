package htmlform

type PipeLabel struct {
	la Label
}

func NewLabel(name string) PipeLabel {
	return PipeLabel{Label{
		Name: name,
	}}
}

func (la PipeLabel) Get() Label {
	return la.la
}

func (la PipeLabel) GetStruct() FormHandlerExt {
	return la.Get()
}

func (la PipeLabel) For(_for string) PipeLabel {
	la.la.For = _for
	return la
}

func (la PipeLabel) Id(id string) PipeLabel {
	la.la.Id = id
	return la
}

func (la PipeLabel) Class(class string) PipeLabel {
	la.la.Class = class
	return la
}
