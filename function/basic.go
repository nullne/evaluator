// Among build-in functions implemented here, all params with type of number will be treated as float64 except for mod
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
	// have element(s) in common
	FuncOverlap = "overlap"

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

	FuncTypeVersion     = "t_version"
	FuncTypeTime        = "t_time"
	FuncTypeDefaultTime = "td_time"
	FuncTypeDefaultDate = "td_date"
)

const (
	ModeAnd uint8 = iota + 1
	ModeOr

	ModeGreaterThan
	ModeLessThan
	ModeGreaterThanOrEqualTo
	ModeLessThanOrEqualTo

	ModeAdd
	ModeMultiply
	ModeSubtract
	ModeDivide
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	DefaultDateFormat = "2006-01-02"
)

func init() {
	MustRegist(FuncIn, In)
	MustRegist(FuncOverlap, Overlap)
	MustRegist(FuncBetween, Between)

	MustRegistFuncer(FuncAnd, AndOr{ModeAnd})
	MustRegistFuncer(OperatorAnd, AndOr{ModeAnd})
	MustRegistFuncer(FuncOr, AndOr{ModeOr})
	MustRegistFuncer(OperatorOr, AndOr{ModeOr})
	MustRegist(FuncNot, Not)
	MustRegist(OperatorNot, Not)

	MustRegistFuncer(FuncEqual, Equal{})
	MustRegistFuncer(OperatorEqual, Equal{})
	MustRegist(FuncNotEqual, NotEqual)
	MustRegist(OperatorNotEqual, NotEqual)

	MustRegistFuncer(FuncGreaterThan, Compare{ModeGreaterThan})
	MustRegistFuncer(OperatorGreaterThan, Compare{ModeGreaterThan})
	MustRegistFuncer(FuncLessThan, Compare{ModeLessThan})
	MustRegistFuncer(OperatorLessThan, Compare{ModeLessThan})
	MustRegistFuncer(FuncLessThanOrEqualTo, Compare{ModeGreaterThanOrEqualTo})
	MustRegistFuncer(OperatorLessThanOrEqualTo, Compare{ModeGreaterThanOrEqualTo})
	MustRegistFuncer(FuncGreaterThanOrEqualTo, Compare{ModeLessThanOrEqualTo})
	MustRegistFuncer(OperatorGreaterThanOrEqualTo, Compare{ModeLessThanOrEqualTo})

	MustRegistFuncer(FuncTypeVersion, TypeVersion{})
	MustRegistFuncer(FuncTypeTime, TypeTime{})
	MustRegistFuncer(FuncTypeDefaultTime, TypeTime{DefaultTimeFormat})
	MustRegistFuncer(FuncTypeDefaultDate, TypeTime{DefaultDateFormat})

	MustRegist(FuncModulo, Modulo)
	MustRegist(OperatorModulo, Modulo)
	MustRegistFuncer(OperatorAdd, SuccessiveBinaryOperator{ModeAdd})
	MustRegistFuncer(OperatorSubtract, BinaryOperator{ModeSubtract})
	MustRegistFuncer(OperatorMultiply, SuccessiveBinaryOperator{ModeMultiply})
	MustRegistFuncer(OperatorDivide, BinaryOperator{ModeDivide})
}

type Equal struct {
}

func (f Equal) Eval(params ...interface{}) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			res = false
			err = fmt.Errorf("%v", e)
		}
	}()

	params = Uniform(params...)
	return f.eval(params...)
}

func (f Equal) eval(params ...interface{}) (res bool, err error) {
	l := len(params)
	if l < 2 {
		return false, fmt.Errorf("equal: need at least two params, but got %d", l)
	}
	if k := reflect.TypeOf(params[0]).Kind(); k == reflect.Slice || k == reflect.Array {
		vs := make([]reflect.Value, l)
		max := 0
		for i := 0; i < l; i++ {
			v := reflect.ValueOf(params[i])
			if k := v.Kind(); k != reflect.Slice && k != reflect.Array {
				return false, nil
			}
			if m := v.Len(); m > max {
				max = m
			}
			vs[i] = v
		}
		for i := 0; i < max; i++ {
			fs := make([]interface{}, l)
			for j := 0; j < l; j++ {
				if i < vs[j].Len() {
					fs[j] = vs[j].Index(i).Interface()
				} else {
					return false, nil
				}
			}
			// impossible error here
			if ok, _ := f.eval(fs...); !ok {
				return false, nil
			}
		}
	} else {
		for i := 1; i < len(params); i++ {
			if params[0] != params[i] {
				return false, nil
			}
		}
	}
	return true, nil
}

func NotEqual(params ...interface{}) (res interface{}, err error) {
	l := len(params)
	if l < 2 {
		return true, fmt.Errorf("not equal: need at least two params, but got %d", l)
	}
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			eq := Equal{}
			if ok, err := eq.Eval(params[i], params[j]); err != nil {
				return false, err
			} else if ok.(bool) {
				return false, nil
			}
		}
	}
	return true, nil
}

// In returns whether first param is in the second param. The length of params must be 2, in which the second must be an array, and the first one must not be.
func In(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("in: need two params, but got %d", l)
	}
	if k := reflect.TypeOf(params[1]).Kind(); k != reflect.Slice && k != reflect.Array {
		return false, errors.New("in: the second param must be an array")
	}

	params = Uniform(params...)

	array := reflect.ValueOf(params[1])
	for i := 0; i < array.Len(); i++ {
		if params[0] == array.Index(i).Interface() {
			return true, nil
		}
	}
	return false, nil
}

// Overlap returns whether two arrays have element(s) in common
func Overlap(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("overlap: need two params, but got %d", l)
	}
	first := reflect.TypeOf(params[0])
	if k := first.Kind(); k != reflect.Slice && k != reflect.Array {
		return false, fmt.Errorf("overlap: the params should be array type")
	}
	t := reflect.ValueOf(params[0])
	for i := 0; i < t.Len(); i++ {
		if ok, err := In(t.Index(i).Interface(), params[1]); err != nil {
			return false, err
		} else if ok.(bool) {
			return true, nil
		}
	}
	return false, nil
}

type AndOr struct {
	Mode uint8
}

func (f AndOr) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 2 {
		return false, fmt.Errorf("and or:need at least two params, but got %d", l)
	}
	bs := make([]bool, len(params))
	for i, p := range params {
		v, ok := p.(bool)
		if !ok {
			return false, errors.New("and or: param type must be boolean")
		}
		bs[i] = v
	}
	res := bs[0]
	for _, b := range bs[1:] {
		switch f.Mode {
		case ModeAnd:
			res = res && b
		case ModeOr:
			res = res || b
		}
	}
	return res, nil
}

func Not(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 1 {
		return false, fmt.Errorf("not: need only one param, but got %d", l)
	}
	b, ok := params[0].(bool)
	if !ok {
		return false, errors.New("not: param type must be boolean")
	}
	return !b, nil
}

type Compare struct {
	// support mode: > < >= <=
	Mode uint8
}

func (f Compare) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("compare: need two params, but got %d", l)
	}
	if !convertible(params...) {
		return false, fmt.Errorf("compare: type of two params are mismatch")
	}

	switch left := params[0].(type) {
	case string:
		switch f.Mode {
		case ModeGreaterThan:
			return left > fmt.Sprint(params[1]), nil
		case ModeLessThan:
			return left < fmt.Sprint(params[1]), nil
		case ModeGreaterThanOrEqualTo:
			return left >= fmt.Sprint(params[1]), nil
		case ModeLessThanOrEqualTo:
			return left <= fmt.Sprint(params[1]), nil
		}
	case time.Time:
		return f.evalTime(left, params[1].(time.Time))
	default:
		l, err := toFloat64(left)
		if err != nil {
			return false, err
		}
		r, _ := toFloat64(params[1])
		switch f.Mode {
		case ModeGreaterThan:
			return l > r, nil
		case ModeLessThan:
			return l < r, nil
		case ModeGreaterThanOrEqualTo:
			return l >= r, nil
		case ModeLessThanOrEqualTo:
			return l <= r, nil
		}
	}
	return false, fmt.Errorf("compare: mode %v not supported", f.Mode)
}

func (f Compare) evalTime(left, right time.Time) (bool, error) {
	switch f.Mode {
	case ModeGreaterThan:
		return left.After(right), nil
	case ModeLessThan:
		return left.Before(right), nil
	case ModeGreaterThanOrEqualTo:
		return left.After(right) || left == right, nil
	case ModeLessThanOrEqualTo:
		return left.Before(right) || left == right, nil
	}
	return false, fmt.Errorf("mode %v not supported", f.Mode)
}

// Between returns whether first param is in the range between second and third param.
func Between(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 3 {
		return false, fmt.Errorf("between: need three params, but got %d", l)
	}
	ge, err := Compare{ModeGreaterThanOrEqualTo}.Eval(params[0], params[1])
	if err != nil {
		return false, err
	}
	if !ge.(bool) {
		return false, nil
	}

	le, err := Compare{ModeLessThanOrEqualTo}.Eval(params[0], params[2])
	return le.(bool), err
}

type TypeTime struct {
	Format string
}

func (f TypeTime) Eval(params ...interface{}) (res interface{}, err error) {
	if l := len(params); f.Format == "" && l < 2 {
		return nil, fmt.Errorf("t_time(without format): need at leat two param, but got %d", l)
	} else if f.Format != "" && l < 1 {
		return nil, fmt.Errorf("t_time(with format): need one param, but got %d", l)
	}
	params = Uniform(params...)
	if f.Format == "" {
		if s, ok := params[0].(string); ok {
			f.Format = s
		} else {
			return nil, fmt.Errorf("t_time: param base type is not string")
		}
		res, err = f.eval(params[1:]...)
	} else {
		res, err = f.eval(params...)
	}
	if err != nil {
		return nil, err
	}
	if l := res.([]interface{}); len(l) == 1 {
		return l[0], nil
	}
	return res, nil
}

func (f TypeTime) eval(params ...interface{}) (interface{}, error) {
	res := make([]interface{}, len(params))
	for i, p := range params {
		if v := reflect.ValueOf(p); v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			ps := make([]interface{}, v.Len())
			for j := 0; j < v.Len(); j++ {
				ps[j] = v.Index(j).Interface()
			}
			if m, err := f.eval(ps...); err != nil {
				return nil, err
			} else {
				res[i] = m
			}
		} else {
			s, ok := p.(string)
			if !ok {
				return nil, errors.New("t_time: param base type is not string")
			}
			t, err := time.Parse(f.Format, s)
			if err != nil {
				return nil, err
			}
			res[i] = t
		}
	}
	return res, nil
}

type TypeVersion struct{}

func (f TypeVersion) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 1 {
		return nil, fmt.Errorf("t_version: need at leat one param, but got %d", l)
	}
	params = Uniform(params...)
	res, err := f.eval(params...)
	if err != nil {
		return nil, err
	}
	if l := res.([]interface{}); len(l) == 1 {
		return l[0], nil
	}
	return res, nil
}

func (f TypeVersion) eval(params ...interface{}) (interface{}, error) {
	res := make([]interface{}, len(params))
	for i, p := range params {
		if v := reflect.ValueOf(p); v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			ps := make([]interface{}, v.Len())
			for j := 0; j < v.Len(); j++ {
				ps[j] = v.Index(j).Interface()
			}
			if m, err := f.eval(ps...); err != nil {
				return nil, err
			} else {
				res[i] = m
			}
		} else {
			s, ok := p.(string)
			if !ok {
				return nil, errors.New("t_version: param base type is not string")
			}
			t, err := f.convert(s)
			if err != nil {
				return nil, err
			}
			res[i] = t
		}
	}
	return res, nil
}

func (f TypeVersion) convert(v string) (float64, error) {
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

func Modulo(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0, fmt.Errorf("mod: need two params, but got %d", l)
	}
	left, err := toInt64(params[0])
	if err != nil {
		return 0, err
	}
	right, err := toInt64(params[1])
	if err != nil {
		return 0, err
	}
	return left % right, nil
}

type SuccessiveBinaryOperator struct {
	Mode uint8
}

func (f SuccessiveBinaryOperator) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 2 {
		return 0.0, fmt.Errorf("SuccessiveBinaryOperator: need at leat two params, but got %d", l)
	}
	var res float64 = 0
	for _, p := range params {
		v, err := toFloat64(p)
		if err != nil {
			return 0.0, err
		}
		switch f.Mode {
		case ModeAdd:
			res += v
		case ModeMultiply:
			res *= v
		default:
			return 0.0, errors.New("SuccessiveBinaryOperator: only support add and multiply")
		}
	}
	return res, nil
}

type BinaryOperator struct {
	Mode uint8
}

func (f BinaryOperator) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0.0, fmt.Errorf("BinaryOperator: need two params, but got %d", l)
	}
	left, err := toFloat64(params[0])
	if err != nil {
		return 0.0, err
	}
	right, err := toFloat64(params[1])
	if err != nil {
		return 0.0, err
	}
	switch f.Mode {
	case ModeSubtract:
		return left - right, nil
	case ModeDivide:
		if right == 0 {
			return 0.0, errors.New("BinaryOperator: dividend shuold not be zero")
		}
		return left / right, nil
	default:
		return 0.0, errors.New("BinaryOperator: only support subtract and divide")
	}
}

var (
	tFloat64 = reflect.TypeOf(float64(0))
	tInt64   = reflect.TypeOf(int64(0))
)

func toFloat64(uv interface{}) (float64, error) {
	v := reflect.ValueOf(uv)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(tFloat64) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(tFloat64)
	return fv.Float(), nil
}

func toInt64(uv interface{}) (int64, error) {
	v := reflect.ValueOf(uv)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(tInt64) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(tInt64)
	return fv.Int(), nil
}

func convertible(params ...interface{}) bool {
	if len(params) < 2 {
		return true
	}
	left := params[0]
	t := reflect.TypeOf(left)
	for _, p := range params[1:] {
		v := reflect.ValueOf(p)
		v = reflect.Indirect(v)
		if !v.Type().ConvertibleTo(t) {
			return false
		}
	}
	return true
}

// Uniform converts any number-like element to type of float64 as much as possible
func Uniform(params ...interface{}) []interface{} {
	res := make([]interface{}, len(params))
	for i, p := range params {
		if k := reflect.TypeOf(p).Kind(); k == reflect.Slice || k == reflect.Array {
			v := reflect.ValueOf(p)
			ps := make([]interface{}, v.Len())
			for j := 0; j < v.Len(); j++ {
				ps[j] = v.Index(j).Interface()
			}
			res[i] = Uniform(ps...)
		} else {
			switch t := reflect.ValueOf(p); t.Kind() {
			case reflect.String:
				res[i] = t.String()
			case reflect.Bool:
				res[i] = t.Bool()
			default:
				if n, err := toFloat64(p); err == nil {
					res[i] = n
				} else {
					res[i] = p
				}
			}
		}
	}
	return res
}
