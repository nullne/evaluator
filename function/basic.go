package function

const (
	FuncIn                 = "in"
	FuncAnd                = "and"
	FuncOr                 = "or"
	FuncNot                = "not"
	FuncEqual              = "eq"
	FuncNotEqual           = "ne"
	OperatorGreaterThan    = "gt"
	FuncLessThan           = "lt"
	FuncGreaterThanOrEqual = "ge"
	FuncLessThanOrEqual    = "le"
	FuncModulo             = "mod"

	OperatorAnd                = "&"
	OperatorOr                 = "|"
	OperatorNot                = "!"
	OperatorEqual              = "="
	OperatorNotEqual           = "!="
	FuncGreaterThan            = ">"
	OperatorLessThan           = "<"
	OperatorGreaterThanOrEqual = ">="
	OperatorLessThanOrEqual    = "<="
	OperatorAdd                = "+"
	OperatorSubtract           = "-"
	OperatorMultiply           = "*"
	OperatorDivide             = "/"
	OperatorModulo             = "%"
)

func init() {
	MustRegist(FuncIn, In{})
	MustRegist(FuncAnd, And{})
	MustRegist(OperatorAnd, And{})
	MustRegist(FuncOr, Or{})
	MustRegist(OperatorOr, Or{})
	MustRegist(FuncNot, Not{})
	MustRegist(OperatorNot, Not{})
	MustRegist(FuncEqual, Eq{})
	MustRegist(OperatorEqual, Eq{})
	MustRegist(FuncNotEqual, NotEq{})
	MustRegist(OperatorNotEqual, NotEq{})
}

type In struct{}

func (f In) Eval(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, ErrParamsInvalid
	}
	array, ok := params[1].([]interface{})
	if !ok {
		return nil, ErrParamsInvalid
	}
	for _, p := range array {
		if params[0] == p {
			return true, nil
		}
	}
	return false, nil
}

func (f In) Helper() string {
	return ""
}

type And struct{}

func (f And) Eval(params ...interface{}) (interface{}, error) {
	return logicAndOr("and", params...)
}

func (f And) Helper() string {
	return ""
}

type Or struct{}

func (f Or) Eval(params ...interface{}) (interface{}, error) {
	return logicAndOr("or", params)
}

func (f Or) Helper() string {
	return ""
}

func logicAndOr(t string, params ...interface{}) (bool, error) {
	if !(len(params) >= 2) {
		return false, ErrParamsInvalid
	}
	bs := make([]bool, len(params))
	for i, p := range params {
		v, ok := p.(bool)
		if !ok {
			return false, ErrParamsInvalid
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

func (f Not) Eval(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return false, ErrParamsInvalid
	}
	b, ok := params[0].(bool)
	if !ok {
		return false, ErrParamsInvalid
	}
	return !b, nil
}

func (f Not) Helper() string {
	return ""
}

type Eq struct{}

func (f Eq) Eval(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return false, ErrParamsInvalid
	}
	return params[0] == params[1], nil
}

func (f Eq) Helper() string {
	return ""
}

type NotEq struct{}

func (f NotEq) Eval(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return false, ErrParamsInvalid
	}
	return params[0] != params[1], nil
}

func (f NotEq) Helper() string {
	return ""
}
