// s-expression
// - treat the empty list, or first element of list is not function (type is string), as data list
package condition

import (
	"errors"

	"github.com/nullne/condition/function"
)

var (
	ErrAmbiguousName = errors.New("ambiguous name used for getting value from function or variable")
	ErrNotFound      = errors.New("neither function not variable found")
)

// VarValue returns the runtime value for variable in expression
type VarValue func(name string) (interface{}, error)

func evaluate(exp sexp, vvf VarValue) (interface{}, error) {
	if l, isList := exp.i.(list); isList {
		if len(l) == 0 {
			return l, nil
		}
		isFunc := false
		if name, ok := l[0].i.(string); ok {
			if _, err := function.Get(name); err == nil {
				isFunc = true
			}
		}

		params := make([]interface{}, 0, len(l))
		tl := l
		if isFunc {
			tl = l[1:]
		}
		for _, p := range tl {
			v, err := evaluate(p, vvf)
			if err != nil {
				return nil, err
			}
			params = append(params, v)
		}
		if isFunc {
			fn, _ := function.Get(l[0].i.(string))
			return fn.Eval(params...)
		} else {
			return append(make([]interface{}, 0, len(params)), params...), nil
		}
	} else {
		if val, ok := exp.i.(string); ok {
			return vvf(val)
		}
		return exp.i, nil
	}
}
