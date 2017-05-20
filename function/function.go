package function

import "errors"

var (
	ErrNotFound       = errors.New("function not implemented")
	ErrParamsInvalid  = errors.New("params are invalid")
	ErrFunctionExists = errors.New("function exists")
)

var (
	functions map[string]Function = make(map[string]Function)
)

func Get(name string) (Function, error) {
	fn, exists := functions[name]
	if !exists {
		return nil, ErrNotFound
	}
	return fn, nil
}

type Function interface {
	Eval(params ...interface{}) (interface{}, error)
	Helper() string
}

func Regist(name string, fn Function) error {
	if _, exist := functions[name]; exist {
		return ErrFunctionExists
	}
	functions[name] = fn
	return nil
}

func MustRegist(name string, fn Function) {
	functions[name] = fn
}
