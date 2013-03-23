package htmlform

type PipeOption struct {
	op Option
}

func NewOption(name string) PipeOption {
	return PipeOption{
		Option{Name: name},
	}
}

func (op PipeOption) Get() *Option {
	return &op.op
}

func (op PipeOption) Value(value string) PipeOption {
	op.op.Value = value
	return op
}

func (op PipeOption) Label(label string) PipeOption {
	op.op.Label = label
	return op
}

func (op PipeOption) Selected() PipeOption {
	op.op.Selected = true
	return op
}
