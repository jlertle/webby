package htmlform

type PipeTextarea struct {
	te Textarea
}

func NewTextarea(name string) PipeTextarea {
	return PipeTextarea{
		Textarea{Name: name},
	}
}

func (te PipeTextarea) Get() *Textarea {
	return &te.te
}

func (te PipeTextarea) Value(value string) PipeTextarea {
	te.te.Value = value
	return te
}

func (te PipeTextarea) Id(id string) PipeTextarea {
	te.te.Id = id
	return te
}

func (te PipeTextarea) Class(class string) PipeTextarea {
	te.te.Class = class
	return te
}

func (te PipeTextarea) MinChar(minchar int) PipeTextarea {
	te.te.MinChar = minchar
	return te
}

func (te PipeTextarea) Mandatory() PipeTextarea {
	return te.MinChar(1)
}

func (te PipeTextarea) MaxChar(maxchar int) PipeTextarea {
	te.te.MaxChar = maxchar
	return te
}

func (te PipeTextarea) Rows(rows int) PipeTextarea {
	te.te.Rows = rows
	return te
}

func (te PipeTextarea) Cols(cols int) PipeTextarea {
	te.te.Cols = cols
	return te
}

func (te PipeTextarea) Extra(extra func(*Validation) error) PipeTextarea {
	te.te.extra = extra
	return te
}
