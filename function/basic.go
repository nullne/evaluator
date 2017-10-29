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

// In returns whether first param is in the second param. The length of params must be 2, in which the second must be an array, and the first one must not be.
func In(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("need two params, but got %d", l)
	}
	array, ok := params[0].([]interface{})
	if ok {
		return false, errors.New("the first param must not be an array")
	}
	array, ok = params[1].([]interface{})
	if !ok {
		return false, errors.New("the second param must be an array")
	}

	isNumber := false
	left, err := toFloat64(params[0])
	if err == nil {
		isNumber = true
	}
	for _, p := range array {
		if isNumber {
			r, err := toFloat64(p)
			if err != nil {
				return false, errors.New("the element in the second param shuold be type of number, same as the first one")
			}
			if left == r {
				return true, nil
			}
		} else {
			if params[0] == p {
				return true, nil
			}
		}
	}
	return false, nil
}

// Overlap returns whether two arrays have element(s) in common
func Overlap(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("need two params, but got %d", l)
	}
	for _, p := range params {
		if _, ok := p.([]interface{}); !ok {
			return false, fmt.Errorf("the params should be array type")
		}
	}
	for _, p := range params[0].([]interface{}) {
		if ok, err := In(p, params[1]); err != nil {
			return false, err
		} else if ok.(bool) {
			return true, nil
		}
	}
	return false, nil
}

// Between returns whether first param is in the range between second and third param.
func Between(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 3 {
		return false, fmt.Errorf("need three params, but got %d", l)
	}
	ge, err := Compare{ModeGreaterThanOrEqualTo}.Eval(params[0], params[1])
	if err != nil {
		return false, err
	}
	if !ge.(bool) {
		return false, nil
	}

	le, err := Compare{ModeLessThanOrEqualTo}.Eval(params[0], params[2])
	if err != nil {
		return false, err
	}
	if !le.(bool) {
		return false, nil
	}
	return true, nil
}

type AndOr struct {
	Mode uint8
}

func (f AndOr) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l < 2 {
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
		return false, fmt.Errorf("need only one param, but got %d", l)
	}
	b, ok := params[0].(bool)
	if !ok {
		return false, errors.New("the type of param must be boolean")
	}
	return !b, nil
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

	if l := len(params); l < 2 {
		return false, fmt.Errorf("need at least two params, but got %d", l)
	}
	if first, ok := params[0].([]interface{}); ok {
		for i := 0; i < len(first); i++ {
			ps := make([]interface{}, len(params))
			for j := 0; j < len(params); j++ {
				ps[j] = params[j].([]interface{})[i]
			}
			res, err := f.Eval(ps...)
			if err != nil {
				return false, err
			}
			if !res.(bool) {
				return false, err
			}
		}
		return true, nil
	} else {
		isNumber := false
		if _, err := toFloat64(params[0]); err == nil {
			isNumber = true
		}
		if isNumber {
			return f.evalFloat64(params...)
		} else {
			for _, p := range params[1:] {
				if p != params[0] {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

func (f Equal) evalFloat64(params ...interface{}) (res bool, err error) {
	left, err := toFloat64(params[0])
	if err != nil {
		return false, err
	}
	for _, p := range params[1:] {
		right, err := toFloat64(p)
		if err != nil {
			return false, err
		}
		if left != right {
			return false, nil
		}
	}
	return true, nil
}

func NotEqual(params ...interface{}) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			res = true
			err = fmt.Errorf("%v", e)
		}
	}()

	if l := len(params); l < 2 {
		return true, fmt.Errorf("need at least two params, but got %d", l)
	}
	pp, err := uniform(params...)
	if err != nil {
		return true, err
	}
	if first, ok := pp[0].([]interface{}); ok {
		list := make([][]interface{}, 0, len(first))
		for i := 0; i < len(first); i++ {
			list = append(list, make([]interface{}, len(pp)))
		}
		for i := 0; i < len(pp); i++ {
			for j := 0; j < len(first); j++ {
				list[j][i] = (pp[i].([]interface{}))[j]
			}
		}
		for _, l := range list {
			r, err := NotEqual(l...)
			if err != nil {
				return true, err
			}
			if !r.(bool) {
				return false, err
			}
		}
	} else {
		for i := 0; i < len(pp); i++ {
			for j := i + 1; j < len(pp); j++ {
				if pp[i] == pp[j] {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

type Compare struct {
	// support mode: > < >= <=
	Mode uint8
}

func (f Compare) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return false, fmt.Errorf("need two params, but got %d", l)
	}
	if !convertible(params...) {
		return false, fmt.Errorf("type of two params are mismatch")
	}
	exists := false
	for _, o := range []uint8{ModeGreaterThan, ModeLessThan, ModeGreaterThanOrEqualTo, ModeLessThanOrEqualTo} {
		if f.Mode == o {
			exists = true
			break
		}
	}
	if !exists {
		return false, fmt.Errorf("mode %v not supported", f.Mode)
	}

	switch left := params[0].(type) {
	case string:
		switch f.Mode {
		case ModeGreaterThan:
			return left > params[1].(string), nil
		case ModeLessThan:
			return left < params[1].(string), nil
		case ModeGreaterThanOrEqualTo:
			return left >= params[1].(string), nil
		case ModeLessThanOrEqualTo:
			return left <= params[1].(string), nil
		}
		return false, nil
	case time.Time:
		return f.evalTime(left, params[1].(time.Time)), nil
	default:
		l, err := toFloat64(left)
		if err != nil {
			return false, err
		}
		r, err := toFloat64(params[1])
		if err != nil {
			return false, err
		}
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
		return false, ErrNotFound
	}
}

func (f Compare) evalTime(left, right time.Time) bool {
	switch f.Mode {
	case ModeGreaterThan:
		return left.After(right)
	case ModeLessThan:
		return left.Before(right)
	case ModeGreaterThanOrEqualTo:
		return left.After(right) || left == right
	case ModeLessThanOrEqualTo:
		return left.Before(right) || left == right
	}
	return false
}

type TypeTime struct {
	Format string
}

func (f TypeTime) Eval(params ...interface{}) (interface{}, error) {
	if f.Format != "" {
		params = append([]interface{}{f.Format}, params...)
	}
	return typeTime{}.eval(params...)
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

func Modulo(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0, fmt.Errorf("need two params, but got %d", l)
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
		return 0.0, fmt.Errorf("need at leat two params, but got %d", l)
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
			return 0.0, errors.New("only support add and multiply")
		}
	}
	return res, nil
}

type BinaryOperator struct {
	Mode uint8
}

func (f BinaryOperator) Eval(params ...interface{}) (interface{}, error) {
	if l := len(params); l != 2 {
		return 0.0, fmt.Errorf("need two params, but got %d", l)
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
			return 0.0, errors.New("dividend shuold not be zero")
		}
		return left / right, nil
	default:
		return 0.0, errors.New("only support subtract and divide")
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

// should handle possible panic when invoked
func uniform(params ...interface{}) ([]interface{}, error) {
	if l := len(params); l < 1 {
		return nil, fmt.Errorf("need at least one params, but got %d", l)
	}
	res := make([]interface{}, 0, len(params))
	if _, ok := params[0].([]interface{}); ok {
		for _, p := range params {
			r, err := uniform(p.([]interface{})...)
			if err != nil {
				return nil, err
			}
			res = append(res, r)
		}
	} else {
		isNumber := false
		if _, err := toFloat64(params[0]); err == nil {
			isNumber = true
		}
		for _, p := range params {
			if isNumber {
				n, err := toFloat64(p)
				if err != nil {
					return nil, err
				}
				res = append(res, n)
			} else {
				res = append(res, p)
			}
		}
	}
	return res, nil
}
