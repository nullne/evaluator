package condition

import "testing"

func TestBasic(t *testing.T) {
	// stmt := `(in gender ("female" "male"))`
	type input struct {
		expr string
		res  interface{}
		err  error
	}
	vvf := func(name string) (interface{}, error) {
		switch name {
		case "gender":
			return "male", nil
		case "age":
			return 18, nil
		case "price":
			return 16.7, nil
		}
		return nil, ErrNotFound
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
		r, err := BoolEval(input.expr, vvf)
		if err != input.err {
			t.Error(err)
		}
		if r != input.res {
			t.Errorf("expression `` %s wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
}
