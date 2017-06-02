package function

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	FuncIn      = "in"
	FuncBetween = "between"

	FuncAnd     = "and"
	FuncOr      = "or"
	FuncNot     = "not"
	OperatorAnd = "&"
	OperatorOr  = "|"
	OperatorNot = "!"

	FuncEqual                    = "eq"
	FuncNotEqual                 = "ne"
	FuncGreaterThan              = "gt"
	FuncLessThan                 = "lt"
	FuncGreaterThanOrEqualTo     = "ge"
	FuncLessThanOrEqualTo        = "le"
	OperatorEqual                = "="
	OperatorNotEqual             = "!="
	OperatorGreaterThan          = ">"
	OperatorLessThan             = "<"
	OperatorGreaterThanOrEqualTo = ">="
	OperatorLessThanOrEqualTo    = "<="

	FuncModulo     = "mod"
	OperatorModulo = "%"

	OperatorAdd      = "+"
	OperatorSubtract = "-"
	OperatorMultiply = "*"
	OperatorDivide   = "/"

	FuncTypeVersion     = "type_version"
	FuncTypeTime        = "type_time"
	FuncTypeDefaultTime = "type_default_time"
	FuncTypeDefaultDate = "type_default_date"
)

func init() {
	MustRegist(FuncIn, In{})
	MustRegist(FuncBetween, Between{})

	MustRegist(FuncAnd, And{})
	MustRegist(OperatorAnd, And{})
	MustRegist(FuncOr, Or{})
	MustRegist(OperatorOr, Or{})
	MustRegist(FuncNot, Not{})
	MustRegist(OperatorNot, Not{})

	MustRegist(FuncEqual, Equal{})
	MustRegist(OperatorEqual, Equal{})
	MustRegist(FuncNotEqual, NotEqual{})
	MustRegist(OperatorNotEqual, NotEqual{})
	MustRegist(FuncGreaterThan, GreaterThan{})
	MustRegist(OperatorGreaterThan, GreaterThan{})
	MustRegist(FuncLessThan, LessThan{})
	MustRegist(OperatorLessThan, LessThan{})
	MustRegist(FuncLessThanOrEqualTo, LessThanOrEqualTo{})
	MustRegist(OperatorLessThanOrEqualTo, LessThanOrEqualTo{})
	MustRegist(FuncGreaterThanOrEqualTo, GreaterThanOrEqualTo{})
	MustRegist(OperatorGreaterThanOrEqualTo, GreaterThanOrEqualTo{})

	MustRegist(FuncTypeVersion, TypeVersion{})
	MustRegist(FuncTypeTime, TypeTime{})
	MustRegist(FuncTypeDefaultTime, TypeDefaultTime{})
	MustRegist(FuncTypeDefaultDate, TypeDefaultDate{})

	MustRegist(FuncModulo, Modulo{})
	MustRegist(OperatorModulo, Modulo{})
	MustRegist(OperatorAdd, Add{})
	MustRegist(OperatorSubtract, Subtract{})
	MustRegist(OperatorMultiply, Multiply{})
	MustRegist(OperatorDivide, Divide{})
}

type In struct{}

// Eval returns whether fisrt param is in the second param. return true only when both type and value are the same, otherwise return false. The length of params must be 2 in which the second should be an array
// In{}.Eval("1", []string{"1", "2"})
func (f In) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("need two params, but got %d", l)
	}
	array, ok := params[1].([]interface{})
	if !ok {
		return false, errors.New("the second param must be an array")
	}
	for _, p := range array {
		if params[0] == p {
			return true, nil
		}
	}
	return false, nil
}

type Between struct{}

// Eval returns whether first param is in the range between second and third param. The type of params can be int, float64, string, and only one of them at one time.
// Between{}.Eval(33, 10, 100)
func (f Between) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 3 {
		return false, fmt.Errorf("need three params, but got %d", l)
	}
	ge, err := GreaterThanOrEqualTo{}.Eval(params[0], params[1])
	if err != nil {
		return false, err
	}
	if !ge.(bool) {
		return false, nil
	}

	le, err := LessThanOrEqualTo{}.Eval(params[0], params[2])
	if err != nil {
		return false, err
	}
	if !le.(bool) {
		return false, nil
	}
	return true, nil
}

type And struct{}

// Eval returns the result of logic AND for params which must be type of boolean
func (f And) Eval(params ...interface{}) (interface{}, error) {
	return logicAndOr("and", params...)
}

type Or struct{}

// Eval returns the result of logic OR for params which must be type of boolean
func (f Or) Eval(params ...interface{}) (interface{}, error) {
	return logicAndOr("or", params...)
}

func logicAndOr(t string, params ...interface{}) (bool, error) {
	if l := len(params); !(l >= 2) {
		return false, fmt.Errorf("need at least two params, but got %d", l)
	}
	bs := make([]bool, len(params))
	for i, p := range params {
		v, ok := p.(bool)
		if !ok {
			return false, errors.New("the type of param must be boolean")
		}
		bs[i] = v
	}
	res := bs[0]
	for _, b := range bs[1:] {
		switch t {
		case "and":
			res = res && b
		case "or":
			res = res || b
		}
	}
	return res, nil
}

type Not struct{}

// Eval returns the result of logic NOT for param which must be type of boolean and the length of params must be 1
func (f Not) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 1 {
		return false, fmt.Errorf("need only one params, but got %d", l)
	}
	b, ok := params[0].(bool)
	if !ok {
		return false, errors.New("the type of param must be boolean")
	}
	return !b, nil
}

type Equal struct{}

// Eval returns true if all params are euqal to each other, otherwise return false
func (f Equal) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 2 {
		return false, fmt.Errorf("need at least two params, but got %d", l)
	}
	for _, p := range params[1:] {
		if p != params[0] {
			return false, nil
		}
	}
	return true, nil
}

type NotEqual struct{}

// Eval returns true if all params are NOT euqal to each other, otherwise return false
func (f NotEqual) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 2 {
		return false, fmt.Errorf("need at least two params, but got %d", l)
	}
	for i := 0; i < len(params); i++ {
		for j := i + 1; j < len(params); j++ {
			if params[i] == params[j] {
				return false, nil
			}
		}
	}
	return true, nil
}

type GreaterThan struct{}

// Eval returns true if first param is greater than the second one, otherwise return false. The length of params must be 2
func (f GreaterThan) Eval(params ...interface{}) (interface{}, error) {
	return (&compare{}).eval(">", params...)
}

type LessThan struct{}

// Eval returns true if first param is less than the second one, otherwise return false. The length of params must be 2
func (f LessThan) Eval(params ...interface{}) (interface{}, error) {
	return (&compare{}).eval("<", params...)
}

type GreaterThanOrEqualTo struct{}

// Eval returns true if first param is greater than or euqal to the second one, otherwise return false. The length of params must be 2
func (f GreaterThanOrEqualTo) Eval(params ...interface{}) (interface{}, error) {
	return (&compare{}).eval(">=", params...)
}

type LessThanOrEqualTo struct{}

// Eval returns true if first param is less than or euqal to the second one, otherwise return false. The length of params must be 2
func (f LessThanOrEqualTo) Eval(params ...interface{}) (interface{}, error) {
	return (&compare{}).eval("<=", params...)
}

type compare struct{}

// compare.eval support operations: > < >= <=
func (f *compare) eval(op string, params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("need two params, but got %d", l)
	}
	if reflect.TypeOf(params[0]) != reflect.TypeOf(params[1]) {
		return false, fmt.Errorf("type of two params are mistatch")
	}
	exists := false
	for _, o := range []string{">", ">=", "<", "<="} {
		if op == o {
			exists = true
		}
	}
	if !exists {
		return false, fmt.Errorf("operand %v not supported", op)
	}

	switch left := params[0].(type) {
	case int:
		return f.evalFloat64(op, float64(left), float64(params[1].(int))), nil
	case float64:
		return f.evalFloat64(op, left, params[1].(float64)), nil
	case string:
		return f.evalString(op, left, params[1].(string)), nil
	case time.Time:
		return f.evalTime(op, left, params[1].(time.Time)), nil

	default:
		return false, ErrNotFound
	}
}

func (f *compare) evalFloat64(op string, left, right float64) bool {
	switch op {
	case ">":
		return left > right
	case "<":
		return left < right
	case ">=":
		return left >= right
	case "<=":
		return left <= right
	}
	return false
}

func (f *compare) evalString(op string, left, right string) bool {
	switch op {
	case ">":
		return left > right
	case "<":
		return left < right
	case ">=":
		return left >= right
	case "<=":
		return left <= right
	}
	return false
}

func (f *compare) evalTime(op string, left, right time.Time) bool {
	switch op {
	case ">":
		return left.After(right)
	case "<":
		return left.Before(right)
	case ">=":
		return left.After(right) || left == right
	case "<=":
		return left.Before(right) || left == right
	}
	return false
}

type TypeTime struct{}

func (f TypeTime) Eval(params ...interface{}) (interface{}, error) {
	return typeTime{}.eval(params...)
}

type TypeDefaultTime struct{}

func (f TypeDefaultTime) Eval(params ...interface{}) (interface{}, error) {
	return typeTime{}.eval(append([]interface{}{"2006-01-02 15:04:05"}, params...)...)
}

type TypeDefaultDate struct{}

func (f TypeDefaultDate) Eval(params ...interface{}) (interface{}, error) {
	return typeTime{}.eval(append([]interface{}{"2006-01-02"}, params...)...)
}

type typeTime struct{}

func (f typeTime) eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return nil, fmt.Errorf("need two param, but got %d", l)
	}
	if list, ok := params[1].([]interface{}); ok {
		res := make([]time.Time, 0, len(list))
		for _, p := range list {
			v, err := f.eval(params[0], p)
			if err != nil {
				return nil, err
			}
			if t, ok := v.(time.Time); ok {
				res = append(res, t)
			} else {
				return nil, ErrIllegalFormat
			}
		}
		return res, nil
	} else {
		v, err := f.convert(params[0], params[1])
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func (f typeTime) convert(ll, vv interface{}) (time.Time, error) {
	l, ok := ll.(string)
	if !ok {
		return time.Time{}, errors.New("basic type of layout value should be string")
	}
	v, ok := vv.(string)
	if !ok {
		return time.Time{}, errors.New("basic type of time value should be string")
	}
	return time.Parse(l, v)
}

type TypeVersion struct{}

func (f TypeVersion) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 1 {
		return false, fmt.Errorf("need only one param, but got %d", l)
	}
	if list, ok := params[0].([]interface{}); ok {
		res := make([]float64, 0, len(list))
		for _, p := range list {
			v, err := f.Eval(p)
			if err != nil {
				return nil, err
			}
			if t, ok := v.(float64); ok {
				res = append(res, t)
			} else {
				return nil, ErrIllegalFormat
			}
		}
		return res, nil
	} else {
		v, err := f.convert(params[0])
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func (f TypeVersion) convert(vv interface{}) (float64, error) {
	v, ok := vv.(string)
	if !ok {
		return 0, errors.New("basic type of version value should be string")
	}
	nums := strings.Split(v, ".")
	if l := len(nums); l < 1 || l > 10 {
		return 0, fmt.Errorf("support at most 10 parts in version")
	}
	var version float64
	e := 5
	for _, num := range nums {
		n, err := strconv.Atoi(num)
		if err != nil {
			return 0, err
		}
		if float64(n) >= math.Pow10(4) {
			return 0, errors.New("each part of version should not greater than 10000")
		}
		version += float64(n) * math.Pow10(4*e)
		e -= 1
	}
	return version, nil
}

type Modulo struct{}

func (f Modulo) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0, fmt.Errorf("need two params, but got %d", l)
	}
	left, ok := params[0].(int)
	if !ok {
		return 0, errors.New("first param should be int")
	}
	right, ok := params[0].(int)
	if !ok {
		return 0, errors.New("second param should be int")
	}
	return left % right, nil
}

type Add struct{}

func (f Add) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 1 {
		return 0.0, fmt.Errorf("need at least one param, but got %d", l)
	}
	var res float64 = 0
	for _, p := range params {
		switch v := p.(type) {
		case int:
			res += float64(v)
		case float64:
			res += v
		default:
			return 0.0, errors.New("the type of the params shuold be int or float")
		}
	}
	return res, nil
}

type Subtract struct{}

func (f Subtract) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0.0, fmt.Errorf("need two params, but got %d", l)
	}
	operators := make([]float64, 2)
	for i, p := range params {
		switch v := p.(type) {
		case int:
			operators[i] = float64(v)
		case float64:
			operators[i] = v
		default:
			return 0.0, errors.New("the type of the params shuold be int or float")
		}
	}
	return operators[0] - operators[1], nil
}

type Multiply struct{}

func (f Multiply) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 1 {
		return 0.0, fmt.Errorf("need at least one param, but got %d", l)
	}
	var res float64 = 0
	for _, p := range params {
		switch v := p.(type) {
		case int:
			res *= float64(v)
		case float64:
			res *= v
		default:
			return 0.0, errors.New("the type of the params shuold be int or float")
		}
	}
	return res, nil
}

type Divide struct{}

func (f Divide) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0.0, fmt.Errorf("need two params, but got %d", l)
	}
	operators := make([]float64, 2)
	for i, p := range params {
		switch v := p.(type) {
		case int:
			operators[i] = float64(v)
		case float64:
			operators[i] = v
		default:
			return 0.0, errors.New("the type of the params shuold be int or float")
		}
	}
	if operators[1] == 0 {
		return 0.0, errors.New("dividend shuold not be zero")
	}
	return operators[0] / operators[1], nil
}
