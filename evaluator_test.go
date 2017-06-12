package evaluator

import "testing"

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
	vvf := MapParams{
		"gender": "male",
		"age":    18,
		"price":  16.7,
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
	// stmt := `(in gender ("female" "male"))`
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
		{`(eq (/ 10 0) 0)`, true},
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
			t.Errorf("expression `` %s wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
}
