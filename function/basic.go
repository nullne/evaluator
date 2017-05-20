package function

type In struct{}

func (f In) Eval(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, ErrParamsInvalid
	}
	array, ok := params[1].([]interface{})
	if !ok {
		return nil, ErrParamsInvalid
	}
	for _, p := range array {
		if params[0] == p {
			return true, nil
		}
	}
	return false, nil
}

func (f In) Help() string {
	return ""
}
