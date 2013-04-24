package webby

type Param map[string]string

func (pa Param) Set(name, value string) {
	pa[name] = value
}

func (pa Param) Get(name string) string {
	return pa[name]
}

func (pa Param) GetInt64(name string) int64 {
	num := int64(0)
	var err error
	num, err = toInt(pa[name])
	if err != nil {
		return 0
	}
	return num
}

func (pa Param) GetInt(name string) int {
	return int(pa.GetInt64(name))
}

func (pa Param) GetInt32(name string) int32 {
	return int32(pa.GetInt64(name))
}

func (pa Param) GetInt16(name string) int16 {
	return int16(pa.GetInt64(name))
}

func (pa Param) GetInt8(name string) int8 {
	return int8(pa.GetInt64(name))
}

func (pa Param) GetUint64(name string) uint64 {
	num := uint64(0)
	var err error
	num, err = toUint(pa[name])
	if err != nil {
		return 0
	}
	return num
}

func (pa Param) GetUint(name string) uint {
	return uint(pa.GetUint64(name))
}

func (pa Param) GetUint32(name string) uint32 {
	return uint32(pa.GetUint64(name))
}

func (pa Param) GetUint16(name string) uint16 {
	return uint16(pa.GetUint64(name))
}

func (pa Param) GetUint8(name string) uint8 {
	return uint8(pa.GetUint64(name))
}

func (pa Param) GetFloat64(name string) float64 {
	num := float64(0)
	var err error
	num, err = toFloat(pa[name])
	if err != nil {
		return float64(0)
	}
	return num
}

func (pa Param) GetFloat32(name string) float32 {
	return float32(pa.GetFloat64(name))
}
