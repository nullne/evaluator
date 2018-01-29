package evaluator

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nullne/evaluator/function"
)

func ExampleExpression() {
	exp, err := New(`(eq gender 'male')`)
	if err != nil {
		log.Fatal(err)
	}
	params := MapParams{"gender": "male"}
	res, err := exp.Eval(params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	res, err = exp.EvalBool(params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	// Output:
	// true
	// true
}

func TestBasic(t *testing.T) {
	params := MapParams{
		"gender": "female",
	}
	res, err := EvalBool(`(in gender ("female" "male"))`, params)
	if err != nil {
		t.Error(err)
	}
	if res != true {
		t.Errorf("incorrect result")
	}
	bres, err := Eval(`(in gender ("female" "male"))`, params)
	if err != nil {
		t.Error(err)
	}
	if bres != true {
		t.Errorf("incorrect result")
	}
}

func TestBasicIncorrect(t *testing.T) {
	params := MapParams{
		"gender": "female",
	}
	expr, err := New(`(+ 1 1)`)
	if err != nil {
		t.Error(err)
	}
	_, err = expr.EvalBool(params)
	if err == nil {
		t.Error("should have errors")
	}
}

func TestComplicated(t *testing.T) {
	appVersion, err := function.TypeVersion{}.Eval("2.7.1")
	if err != nil {
		t.Fatal(err)
	}
	params := MapParams{
		"gender":      "female",
		"age":         55,
		"app_version": appVersion,
		"region":      []int{1, 2, 3},
	}
	expr, err := New(`
(or
	(and
	(between age 18 80)
	(eq gender "male")
	(between app_version (t_version "2.7.1") (t_version "2.9.1"))
	)
	(overlap region (2890 3780))
 )`)
	if err != nil {
		t.Error(err)
	}
	_, err = expr.EvalBool(params)
	if err != nil {
		t.Error(err)
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
		{`(in gender () )`, false},
		{`(not (in gender ("female" "male")))`, false},
		{`(! (in gender ("female" "male")))`, false},
		{`(ge (t_version "2.1.1") (t_version "2.1.1"))`, true},
		{`(gt (t_version "2.1.1") (t_version "2.1.1"))`, false},
		{`(between (t_version "2.1.1") (t_version "2.1.1") (t_version "2.1.1"))`, true},
		{`(between (t_version "2.1.1.9999") (t_version "2.1.1") (t_version "2.1.2"))`, true},
		{`(between (mod age 5) 1 3)`, true},
		{`(between (td_time now1) (td_time now1) (td_time now2))`, true},
		{`(ge (t_version "2.8.1") (t_version "2.9.3"))`, false},
		{`(ge (t_version "2.9.1") (t_version "2.8.3"))`, true},

		// overlap
		{`(overlap (1 2 3) (4 5 6))`, false},
		{`(overlap () ())`, false},
		{`(overlap (1 2 3) (4 3 2))`, true},
		{`(overlap ("1" "2" "3") ("4" "3" "2"))`, true},
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

		r, err = Eval(input.expr, vvf)
		if err != nil {
			t.Error(err)
		}
		if r != input.res {
			t.Errorf("expression `%s` wanna: %+v, got: %+v", input.expr, input.res, r)
		}
	}
}

func TestIncorrect(t *testing.T) {
	type input struct {
		expr string
		// expression correctness
		b bool
	}
	vvf := MapParams{
		"gender": "male",
		"age":    18,
		"price":  16.7,
	}
	inputs := []input{
		{`eq (mod age 5) 3.0`, false},
		{`(eq (+ money 5) 15)`, true},
	}
	for _, input := range inputs {
		_, err := New(input.expr)
		if !input.b && err == nil {
			t.Error("should have errors")
		}

		_, err = Eval(input.expr, vvf)
		if err == nil {
			t.Errorf("input: %v, should have errors", input.expr)
		}

		_, err = EvalBool(input.expr, vvf)
		if err == nil {
			t.Errorf("input: %v, should have errors", input.expr)
		}
	}
}

func TestAdvancedFunc(t *testing.T) {
	invoker := func(params ...interface{}) (interface{}, error) {
		fn := params[0].(function.Func)
		return fn(params[1:]...)
	}

	if err := function.Regist("invoke", invoker); err != nil {
		t.Error(err)
	}

	exp := `(invoke + 1 1)`
	e, err := New(exp)
	if err != nil {
		t.Error(err)
	}
	r, err := e.Eval(nil)
	if err != nil {
		t.Error(err)
	}
	if r.(float64) != 2 {
		t.Errorf("expression `%s` wanna: %+v, got: %+v", exp, 2, r)
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

func BenchmarkEqualString(b *testing.B) {
	expr, err := New(`(eq "one" "one" "three")`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkInString(b *testing.B) {
	expr, err := New(`(in "2.7.2" ("2.7.1" "2.7.4" "2.7.5"))`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkBetweenInt(b *testing.B) {
	expr, err := New(`(between 100 10 200)`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkBetweenTime(b *testing.B) {
	expr, err := New(`(between (td_time "2017-07-09 12:00:00") (td_time "2017-07-02 12:00:00") (td_time "2017-07-19 12:00:00"))`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkOverlapInt(b *testing.B) {
	expr, err := New(`(overlap (1 2 3) (4 5 6))`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.EvalBool(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTypeTime(b *testing.B) {
	expr, err := New(`(td_time "2017-09-09 12:00:00")`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.Eval(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTypeVersion(b *testing.B) {
	expr, err := New(`(t_version "2.8.9.1")`)
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		_, err := expr.Eval(nil)
		if err != nil {
			b.Error(err)
		}
	}
}
