// s-expression
// - treat the empty list, or first element of list is not function (type is string), as data list
package condition

import "errors"

var (
	ErrNotFound      = errors.New("neither function not variable found")
	ErrInvalidResult = errors.New("invalid result")
)

// VarValue returns the runtime value for variable in expression
type VarValue func(name string) (interface{}, error)

func Eval(expr string, fn VarValue) (interface{}, error) {
	return eval(expr, fn)
}

func BoolEval(expr string, fn VarValue) (bool, error) {
	r, err := eval(expr, fn)
	if err != nil {
		return false, err
	}
	b, ok := r.(bool)
	if !ok {
		return false, ErrInvalidResult
	}
	return b, nil
}

func eval(expr string, fn VarValue) (interface{}, error) {
	exp, err := parse(expr)
	if err != nil {
		return false, err
	}
	return exp.evaluate(fn)
}
