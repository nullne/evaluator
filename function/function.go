package function

import "errors"

var (
	ErrNotFound      = errors.New("function not implemented")
	ErrParamsInvalid = errors.New("params are invalid")
)

var (
	functions map[string]Function = make(map[string]Function)
)

func Get(name string) (Function, error) {
	return In{}, nil
}

type Function interface {
	Eval(params ...interface{}) (interface{}, error)
	Help() string
}
