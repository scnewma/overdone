package optional

type optional struct {
	value interface{}
}

func Of(value interface{}) optional {
	return optional{value}
}

func (o optional) OrElse(value interface{}) interface{} {
	if o.value != nil {
		return o.value
	} else {
		return value
	}
}

func (o optional) IsPresent() bool {
	return o.value != nil
}
