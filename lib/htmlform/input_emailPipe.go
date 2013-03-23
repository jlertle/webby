package htmlform

type PipeInputEmail struct {
	em InputEmail
}

func NewInputEmail(name string) PipeInputEmail {
	return PipeInputEmail{InputEmail{Name: name}}
}

func (em PipeInputEmail) Get() *InputEmail {
	return &em.em
}

func (em PipeInputEmail) Value(value string) PipeInputEmail {
	em.em.Value = value
	return em
}

func (em PipeInputEmail) Id(id string) PipeInputEmail {
	em.em.Id = id
	return em
}

func (em PipeInputEmail) Class(class string) PipeInputEmail {
	em.em.Class = class
	return em
}

func (em PipeInputEmail) MustMatch(match, errMsg string) PipeInputEmail {
	em.em.MustMatch = match
	em.em.MustMatchErr = errMsg
	return em
}
