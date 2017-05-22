// s-expression
// - treat the empty list, or first element of list is not function (type is string), as data list
package evaluator

import "errors"

var (
	ErrNotFound      = errors.New("neither function not variable found")
	ErrInvalidResult = errors.New("invalid result type")
)

type Expression struct {
	exp sexp
}

func New(expr string) (Expression, error) {
	exp, err := parse(expr)
	if err != nil {
		return Expression{}, err
	}
	return Expression{
		exp: exp,
	}, nil
}

func (e Expression) Eval(params map[string]interface{}) (interface{}, error) {
	return e.exp.evaluate(toParamsFunc(params))
}

func (e Expression) EvalWithParamsFunc(pf ParamsFunc) (interface{}, error) {
	return e.exp.evaluate(pf)
}

func (e Expression) EvalBool(params map[string]interface{}) (bool, error) {
	return e.evalBool(toParamsFunc(params))
}

func (e Expression) EvalBoolWithParamsFunc(pf ParamsFunc) (bool, error) {
	return e.evalBool(pf)
}

func (e Expression) evalBool(pf ParamsFunc) (bool, error) {
	r, err := e.exp.evaluate(pf)
	if err != nil {
		return false, err
	}
	b, ok := r.(bool)
	if !ok {
		return false, ErrInvalidResult
	}
	return b, nil
}

// ParamsFunc returns the runtime value for variable in expression
type ParamsFunc func(name string) (interface{}, error)

func toParamsFunc(params map[string]interface{}) ParamsFunc {
	return func(name string) (interface{}, error) {
		v, ok := params[name]
		if !ok {
			return nil, ErrNotFound
		}
		return v, nil
	}
}

func Eval(expr string, params map[string]interface{}) (interface{}, error) {
	e, err := New(expr)
	if err != nil {
		return nil, err
	}
	return e.Eval(params)
}

func EvalBool(expr string, params map[string]interface{}) (bool, error) {
	e, err := New(expr)
	if err != nil {
		return false, err
	}
	return e.EvalBool(params)
}

func EvalWithParamsFunc(expr string, pf ParamsFunc) (interface{}, error) {
	e, err := New(expr)
	if err != nil {
		return nil, err
	}
	return e.EvalWithParamsFunc(pf)
}

func EvalBoolWithParamsFunc(expr string, pf ParamsFunc) (bool, error) {
	e, err := New(expr)
	if err != nil {
		return false, err
	}
	return e.EvalBoolWithParamsFunc(pf)
}
