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
	// FuncIn is the function/operator keyword in
	FuncIn = "in"
	// FuncBetween is the function/operator keyword between
	FuncBetween = "between"
	// FuncOverlap is the function/operator keyword overlap
	FuncOverlap = "overlap"

	// FuncAnd is the function/operator keyword and
	FuncAnd = "and"
	// FuncOr is the function/operator keyword or
	FuncOr = "or"
	// FuncNot is the function/operator keyword not
	FuncNot = "not"
	// OperatorAnd is the function/operator keyword &
	OperatorAnd = "&"
	// OperatorOr is the function/operator keyword |
	OperatorOr = "|"
	// OperatorNot is the function/operator keyword !
	OperatorNot = "!"

	// FuncEqual is the function/operator keyword eq
	FuncEqual = "eq"
	// FuncNotEqual is the function/operator keyword ne
	FuncNotEqual = "ne"
	// FuncGreaterThan is the function/operator keyword gt
	FuncGreaterThan = "gt"
	// FuncLessThan is the function/operator keyword lt
	FuncLessThan = "lt"
	// FuncGreaterThanOrEqualTo is the function/operator keyword ge
	FuncGreaterThanOrEqualTo = "ge"
	// FuncLessThanOrEqualTo is the function/operator keyword le
	FuncLessThanOrEqualTo = "le"
	// OperatorEqual is the function/operator keyword =
	OperatorEqual = "="
	// OperatorNotEqual is the function/operator keyword !=
	OperatorNotEqual = "!="
	// OperatorGreaterThan is the function/operator keyword >
	OperatorGreaterThan = ">"
	// OperatorLessThan is the function/operator keyword <
	OperatorLessThan = "<"
	// OperatorGreaterThanOrEqualTo is the function/operator keyword >=
	OperatorGreaterThanOrEqualTo = ">="
	// OperatorLessThanOrEqualTo is the function/operator keyword <=
	OperatorLessThanOrEqualTo = "<="

	// FuncModulo is the function/operator keyword mod
	FuncModulo = "mod"
	// OperatorModulo is the function/operator keyword %
	OperatorModulo = "%"

	// OperatorAdd is the function/operator keyword +
	OperatorAdd = "+"
	// OperatorSubtract is the function/operator keyword -
	OperatorSubtract = "-"
	// OperatorMultiply is the function/operator keyword *
	OperatorMultiply = "*"
	// OperatorDivide is the function/operator keyword /
	OperatorDivide = "/"

	// FuncTypeVersion is the function/operator keyword t_version
	FuncTypeVersion = "t_version"
	// FuncTypeTime is the function/operator keyword t_time
	FuncTypeTime = "t_time"
	// FuncTypeDefaultTime is the function/operator keyword td_time
	FuncTypeDefaultTime = "td_time"
	// FuncTypeDefaultDate is the function/operator keyword td_date
	FuncTypeDefaultDate = "td_date"
)

const (
	// ModeAnd is the mode add for function AndOr
	ModeAnd uint8 = iota + 1
	// ModeOr is the mode or for function AndOr
	ModeOr

	// ModeGreaterThan is the mod > for function Compare
	ModeGreaterThan
	// ModeLessThan is the mod < for function Compare
	ModeLessThan
	// ModeGreaterThanOrEqualTo is the mod >= for function Compare
	ModeGreaterThanOrEqualTo
	// ModeLessThanOrEqualTo is the mod <= for function Compare
	ModeLessThanOrEqualTo

	// ModeAdd is the mode + for function SuccessiveBinaryOperator
	ModeAdd
	// ModeMultiply is the mode * for function SuccessiveBinaryOperator
	ModeMultiply
	// ModeSubtract is the mode - for function BinaryOperator
	ModeSubtract
	// ModeDivide is the mode / for function BinaryOperator
	ModeDivide
)

const (
	// DefaultTimeFormat is the default time format
	DefaultTimeFormat = "2006-01-02 15:04:05"
	// DefaultDateFormat is the default date format
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
	MustRegistFuncer(FuncLessThanOrEqualTo, Compare{ModeLessThanOrEqualTo})
	MustRegistFuncer(OperatorLessThanOrEqualTo, Compare{ModeLessThanOrEqualTo})
	MustRegistFuncer(FuncGreaterThanOrEqualTo, Compare{ModeGreaterThanOrEqualTo})
	MustRegistFuncer(OperatorGreaterThanOrEqualTo, Compare{ModeGreaterThanOrEqualTo})

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

// Equal returns whether the input params are equal to each other, array type is supported too
type Equal struct{}

// Eval implements the interface Funcer
func (f Equal) Eval(params ...interface{}) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			res = false
			err = fmt.Errorf("euqal: %v", e)
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

// NotEqual returns whether the input params are not equal with each other, array type is supported too
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

// In returns whether first param is in the second param(must be array type). The length of params must be 2
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

// Overlap returns whether two arrays have element(s) in common. The length of params must be 2 and type must be array.
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

// AndOr implements logic operator "and" "or"
type AndOr struct {
	Mode uint8
}

// Eval implements the interface Funcer
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

// Not implements logic operator "not"
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

// Compare compare the two inputs with the given mode.
type Compare struct {
	// support mode: > < >= <=
	Mode uint8
}

// Eval implements the interface Funcer
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
			return false, fmt.Errorf("compare: %+v", err)
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

// Between returns whether first param is in the range between second and third param. The input params must be comparable
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

// TypeTime converts given string to the time.Time on the format of Format
type TypeTime struct {
	Format string
}

// Eval implements the interface Funcer
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
			m, err := f.eval(ps...)
			if err != nil {
				return nil, err
			}
			res[i] = m
		} else {
			s, ok := p.(string)
			if !ok {
				return nil, errors.New("t_time: param base type is not string")
			}
			t, err := time.Parse(f.Format, s)
			if err != nil {
				return nil, fmt.Errorf("t_time: %+v", err)
			}
			res[i] = t
		}
	}
	return res, nil
}

// TypeVersion converts version string to a comparable number
type TypeVersion struct{}

// Eval implements the interface Funcer
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
			m, err := f.eval(ps...)
			if err != nil {
				return nil, err
			}
			res[i] = m
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
		return 0, fmt.Errorf("t_version: support at most 10 parts in version")
	}
	var version float64
	e := 5
	for _, num := range nums {
		n, err := strconv.Atoi(num)
		if err != nil {
			return 0, fmt.Errorf("t_version: %+v", err)
		}
		if float64(n) >= math.Pow10(4) {
			return 0, errors.New("t_version: each part of version should not greater than 10000")
		}
		version += float64(n) * math.Pow10(4*e)
		e--
	}
	return version, nil
}

// Modulo implements mod operator in go
func Modulo(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0, fmt.Errorf("mod: need two params, but got %d", l)
	}
	left, err := toInt64(params[0])
	if err != nil {
		return 0, fmt.Errorf("mod: %+v", err)
	}
	right, err := toInt64(params[1])
	if err != nil {
		return 0, fmt.Errorf("mod: %+v", err)
	}
	return left % right, nil
}

// SuccessiveBinaryOperator implements successive plus or multiply
type SuccessiveBinaryOperator struct {
	Mode uint8
}

// Eval implements the interface Funcer
func (f SuccessiveBinaryOperator) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 2 {
		return 0.0, fmt.Errorf("SuccessiveBinaryOperator: need at leat two params, but got %d", l)
	}
	var res float64
	for _, p := range params {
		v, err := toFloat64(p)
		if err != nil {
			return 0.0, fmt.Errorf("SuccessiveBinaryOperator: %+v", err)
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

// BinaryOperator implements minus and divide
type BinaryOperator struct {
	Mode uint8
}

// Eval implements the interface Funcer
func (f BinaryOperator) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0.0, fmt.Errorf("BinaryOperator: need two params, but got %d", l)
	}
	left, err := toFloat64(params[0])
	if err != nil {
		return 0.0, fmt.Errorf("BinaryOperator: %+v", err)
	}
	right, err := toFloat64(params[1])
	if err != nil {
		return 0.0, fmt.Errorf("BinaryOperator: %+v", err)
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
