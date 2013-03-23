package htmlform

type PipeInputFile struct {
	fi InputFile
}

func NewInputFile(name string) PipeInputFile {
	return PipeInputFile{InputFile{Name: name}}
}

func (fi PipeInputFile) Get() *InputFile {
	return &fi.fi
}

func (fi PipeInputFile) Id(id string) PipeInputFile {
	fi.fi.Id = id
	return fi
}

func (fi PipeInputFile) Class(class string) PipeInputFile {
	fi.fi.Class = class
	return fi
}

func (fi PipeInputFile) Mime(mimes ...string) PipeInputFile {
	fi.fi.Mime = mimes
	return fi
}

func (fi PipeInputFile) Mandatory() PipeInputFile {
	fi.fi.Mandatory = true
	return fi
}

func (fi PipeInputFile) Size(nbyte int64, errMsg string) PipeInputFile {
	fi.fi.Size = nbyte
	fi.fi.SizeErr = errMsg
	return fi
}
