// Package evaluator evaluates an expression in the form of s-expression
package evaluator

import "errors"

var (
	// ErrNotFound means the unknow string within the expression cannot be Get from neither functions or params
	ErrNotFound = errors.New("neither function not variable found")
	// ErrInvalidResult means the invalid result type expected with the real output
	ErrInvalidResult = errors.New("invalid result type")
)

// Expression stands for an expression which can be evaluated by passing required params
type Expression struct {
	exp sexp
}

// New will return a Expression by parsing the given expression string
func New(expr string) (Expression, error) {
	exp, err := parse(expr)
	if err != nil {
		return Expression{}, err
	}
	return Expression{
		exp: exp,
	}, nil
}

// Eval evaluates the Expression with params and return the real value in the type of interface
func (e Expression) Eval(params Params) (interface{}, error) {
	return e.exp.evaluate(params)
}

// EvalBool invokes method Eval and does boolean type assertion, return ErrInvalidResult if the type of result is not boolean
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

// MapParams is a simple map implementation of Params interface
type MapParams map[string]interface{}

// Get is the only method required by Params interface
func (p MapParams) Get(name string) (interface{}, error) {
	v, ok := p[name]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// Params defines a Get method which gets required param for the expression
type Params interface {
	Get(name string) (interface{}, error)
}

// Eval is a handy encapsulation to parse the expression and evaluate it
func Eval(expr string, params Params) (interface{}, error) {
	e, err := New(expr)
	if err != nil {
		return nil, err
	}
	return e.Eval(params)
}

// EvalBool is same as Eval but return a boolean result instead of interface type
func EvalBool(expr string, params Params) (bool, error) {
	e, err := New(expr)
	if err != nil {
		return false, err
	}
	return e.EvalBool(params)
}
