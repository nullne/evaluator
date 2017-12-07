// Package function provides basic functions which implement Funcer interface
package function

import "errors"

var (
	// ErrNotFound means the function is not implemented
	ErrNotFound = errors.New("function not implemented")
	// ErrParamsInvalid means the passed params are invalid
	ErrParamsInvalid = errors.New("params are invalid")
	// ErrFunctionExists means the function exists when registering
	ErrFunctionExists = errors.New("function exists")
	// ErrIllegalFormat means the data format is not legal
	ErrIllegalFormat = errors.New("illegal data format")
)

var (
	functions = make(map[string]Func)
)

// Get get a registered function by name
func Get(name string) (Func, error) {
	fn, exists := functions[name]
	if !exists {
		return nil, ErrNotFound
	}
	return fn, nil
}

// Funcer is the function interface which will be used to evaluate expressionn
type Funcer interface {
	Eval(params ...interface{}) (interface{}, error)
}

// Func is a handy function type
type Func func(params ...interface{}) (interface{}, error)

// RegistFuncer regists fn with type Funcer with name of name
func RegistFuncer(name string, fn Funcer) error {
	if _, exist := functions[name]; exist {
		return ErrFunctionExists
	}
	functions[name] = fn.Eval
	return nil
}

// MustRegistFuncer is same as RegistFuncer but may overide if function with name existed
func MustRegistFuncer(name string, fn Funcer) {
	functions[name] = fn.Eval
}

// Regist regists fn with type Func with name of name
func Regist(name string, fn Func) error {
	if _, exist := functions[name]; exist {
		return ErrFunctionExists
	}
	functions[name] = fn
	return nil
}

// Registered returns all registered functions or operators
func Registered() []string {
	ss := make([]string, 0, len(functions))
	for k := range functions {
		ss = append(ss, k)
	}
	return ss
}

// MustRegist is same as Regist but may overide if function with name existed
func MustRegist(name string, fn Func) {
	functions[name] = fn
}
