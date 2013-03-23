package htmlform

type PipeSelect struct {
	se Select
}

func NewSelect(name string) PipeSelect {
	return PipeSelect{
		Select{Name: name},
	}
}

func (se PipeSelect) Get() *Select {
	return &se.se
}

func (se PipeSelect) Id(id string) PipeSelect {
	se.se.Id = id
	return se
}

func (se PipeSelect) Class(class string) PipeSelect {
	se.se.Class = class
	return se
}

func (se PipeSelect) Options(options ...*Option) PipeSelect {
	se.se.Options = options
	return se
}

func (se PipeSelect) Mandatory() PipeSelect {
	se.se.Mandatory = true
	return se
}
