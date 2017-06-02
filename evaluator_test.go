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

func TestFuncs(t *testing.T) {
	// stmt := `(in gender ("female" "male"))`
	type input struct {
		expr string
		res  interface{}
		err  error
	}
	vvf := MapParams{
		"gender": "male",
		"age":    18,
		"price":  16.7,
	}
	inputs := []input{
		{
			`(in gender ("female" "male"))`,
			true,
			nil,
		},
		{
			`(in age (16 17))`,
			false,
			nil,
		},
		{
			`(and (in age (16 17)) (in gender ("female" "male")) )`,
			false,
			nil,
		},
		{
			`(in price (16.7 17))`,
			true,
			nil,
		},
	}
	for _, input := range inputs {
		r, err := EvalBool(input.expr, vvf)
		if err != input.err {
			t.Error(err)
		}
		if r != input.res {
			t.Errorf("expression `` %s wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
	for _, input := range inputs {
		e, err := New(input.expr)
		if err != nil {
			t.Error(err)
		}
		r, err := e.EvalBool(vvf)
		if err != input.err {
			t.Error(err)
		}
		if r != input.res {
			t.Errorf("expression `` %s wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
}
