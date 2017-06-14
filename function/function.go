package function

import "errors"

var (
	ErrNotFound       = errors.New("function not implemented")
	ErrParamsInvalid  = errors.New("params are invalid")
	ErrFunctionExists = errors.New("function exists")
	ErrIllegalFormat  = errors.New("illegal data format")
)

var (
	functions map[string]Func = make(map[string]Func)
)

func Get(name string) (Func, error) {
	fn, exists := functions[name]
	if !exists {
		return nil, ErrNotFound
	}
	return fn, nil
}

type Funcer interface {
	Eval(params ...interface{}) (interface{}, error)
}

type Func func(params ...interface{}) (interface{}, error)

func RegistFuncer(name string, fn Funcer) error {
	if _, exist := functions[name]; exist {
		return ErrFunctionExists
	}
	functions[name] = fn.Eval
	return nil
}

func MustRegistFuncer(name string, fn Funcer) {
	functions[name] = fn.Eval
}

func Regist(name string, fn Func) error {
	if _, exist := functions[name]; exist {
		return ErrFunctionExists
	}
	functions[name] = fn
	return nil
}

func MustRegist(name string, fn Func) {
	functions[name] = fn
}
