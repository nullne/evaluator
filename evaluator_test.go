package evaluator

import (
	"errors"
	"testing"
	"time"

	"github.com/nullne/evaluator/function"
)

func TestBasic(t *testing.T) {
	params := MapParams{
		"gender": "female",
	}
	res, err := EvalBool(`(in gender ("female" "male"))`, params)
	if err != nil {
		t.Error(err)
	}
	if res != true {
		t.Errorf("fuck")
	}
}

func TestCorrectBooleanFuncs(t *testing.T) {
	// stmt := `(in gender ("female" "male"))`
	type input struct {
		expr string
		res  interface{}
	}
	now1 := time.Now().Format("2006-01-02 15:04:05")
	now2 := time.Now().Format("2006-01-02 15:04:05")
	vvf := MapParams{
		"gender": "male",
		"age":    18,
		"price":  16.7,
		"now1":   now1,
		"now2":   now2,
	}
	inputs := []input{
		{`(in gender ("female" "male"))`, true},
		{`(not (in gender ("female" "male")))`, false},
		{`(! (in gender ("female" "male")))`, false},
		{`(ge (type_version "2.1.1") (type_version "2.1.1"))`, true},
		{`(gt (type_version "2.1.1") (type_version "2.1.1"))`, false},
		{`(between (type_version "2.1.1") (type_version "2.1.1") (type_version "2.1.1"))`, true},
		{`(between (type_version "2.1.1.9999") (type_version "2.1.1") (type_version "2.1.2"))`, true},
		{`(between (mod age 5) 1 3)`, true},
		{`(between (type_default_time now1) (type_default_time now1) (type_default_time now2))`, true},
	}
	for _, input := range inputs {
		e, err := New(input.expr)
		if err != nil {
			t.Error(err)
		}
		r, err := e.EvalBool(vvf)
		if err != nil {
			t.Error(err)
		}
		if r != input.res {
			t.Errorf("expression `` %s wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
}

func TestCorrectFuncs(t *testing.T) {
	type input struct {
		expr string
		res  interface{}
	}
	vvf := MapParams{
		"gender": "male",
		"age":    18,
		"price":  16.7,
	}
	inputs := []input{
		{`(eq (mod age 5) 3.0)`, true},
		{`(eq (+ 10 5) 15)`, true},
		{`(eq (/ 10 5) 2)`, true},
	}
	for _, input := range inputs {
		e, err := New(input.expr)
		if err != nil {
			t.Error(err)
		}
		r, err := e.Eval(vvf)
		if err != nil {
			t.Error(err)
		}
		if r != input.res {
			t.Errorf("expression `%s` wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
}

type age struct{}

func (f age) Eval(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("only one params accepted")
	}
	birth, ok := params[0].(string)
	if !ok {
		return nil, errors.New("birth format need to be string")
	}
	r, err := time.Parse("2006-01-02", birth)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a := r.Year() - now.Year()
	if r.Month() < now.Month() {
		a--
	} else if r.Month() == now.Month() {
		if r.Day() < now.Day() {
			a--
		}
	}
	return a, nil
}

func TestDIVFunc(t *testing.T) {
	if err := function.Regist("age", age{}); err != nil {
		t.Error(err)
	}

	exp := `(not (between (age birthdate) 18 20))`
	vvf := MapParams{
		"birthdate": "2018-02-01",
	}
	e, err := New(exp)
	if err != nil {
		t.Error(err)
	}
	r, err := e.Eval(vvf)
	if err != nil {
		t.Error(err)
	}
	if r != true {
		t.Errorf("expression `%s` wanna: %+v, got: %+v", exp, true, r)
	}
}
