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

func (e Expression) Eval(params Params) (interface{}, error) {
	return e.exp.evaluate(params)
}

func (e Expression) EvalBool(params Params) (bool, error) {
	r, err := e.exp.evaluate(params)
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
type MapParams map[string]interface{}

func (p MapParams) Get(name string) (interface{}, error) {
	v, ok := p[name]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

type Params interface {
	Get(name string) (interface{}, error)
}

func Eval(expr string, params Params) (interface{}, error) {
	e, err := New(expr)
	if err != nil {
		return nil, err
	}
	return e.Eval(params)
}

func EvalBool(expr string, params Params) (bool, error) {
	e, err := New(expr)
	if err != nil {
		return false, err
	}
	return e.EvalBool(params)
}
