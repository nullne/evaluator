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
		{`(ge (t_version "2.1.1") (t_version "2.1.1"))`, true},
		{`(gt (t_version "2.1.1") (t_version "2.1.1"))`, false},
		{`(between (t_version "2.1.1") (t_version "2.1.1") (t_version "2.1.1"))`, true},
		{`(between (t_version "2.1.1.9999") (t_version "2.1.1") (t_version "2.1.2"))`, true},
		{`(between (mod age 5) 1 3)`, true},
		{`(between (td_time now1) (td_time now1) (td_time now2))`, true},
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

func TestDIVFunc(t *testing.T) {
	age := func(params ...interface{}) (interface{}, error) {
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

	if err := function.Regist("age", age); err != nil {
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

func BenchmarkEvalTimeBetweenReuseExpression(b *testing.B) {
	expr, err := New(`(between (td_date "2017-01-05") (td_date "2017-01-01") (td_date "2017-01-10"))`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		r, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
		if r != true {
			b.Error("wrong result")
		}
	}
}

func BenchmarkEvalTimeBetween(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r, err := EvalBool(`(between (td_date "2017-01-05") (td_date "2017-01-01") (td_date "2017-01-10"))`, nil)
		if err != nil {
			b.Error(err)
		}
		if r != true {
			b.Error("wrong result")
		}
	}
}

func BenchmarkEvalEqualReuseExpression(b *testing.B) {
	expr, err := New(`(eq "male" "female")`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		r, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
		if r != false {
			b.Error("wrong result")
		}
	}
}

func BenchmarkEvalEqual(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r, err := EvalBool(`(eq "male" "female")`, nil)
		if err != nil {
			b.Error(err)
		}
		if r != false {
			b.Error("wrong result")
		}
	}
}
