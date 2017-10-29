package function

import (
	"testing"
	"time"
)

var (
	now, now1, now2  time.Time
	format1, format2 string
)

func init() {
	now := time.Now().UTC()
	format1 = "2006-01-02 15:04:05"
	format2 = "2006-01-02"
	now1, _ = time.Parse(format1, now.Format(format1))
	now2, _ = time.Parse(format2, now.Format(format2))
}

type res struct {
	params []interface{}
	result interface{}
	err    bool
}

func TestIn(t *testing.T) {
	inputs := []res{
		{[]interface{}{100, []interface{}{100, 200, 300}}, true, false},
		{[]interface{}{now, []interface{}{now, time.Now()}}, true, false},
		{[]interface{}{100.00000000000000000000000000000000001, []interface{}{100.00000000000000000000000000000000001, 200.0, 300.0}}, true, false},
		{[]interface{}{"100", []interface{}{"100", "200", "300"}}, true, false},
		{[]interface{}{100.0, []interface{}{100, 200, 300}}, true, false},
		{[]interface{}{true, []interface{}{false, true}}, true, false},

		{[]interface{}{"100", []interface{}{100, 200, 300}}, false, false},
		{[]interface{}{now, []interface{}{time.Now(), time.Now()}}, false, false},

		{[]interface{}{true, 100, []interface{}{false, true}}, false, true},
		{[]interface{}{[]interface{}{true}, false}, false, true},
		{[]interface{}{true, false}, false, true},
		{[]interface{}{[]interface{}{false, true}}, false, true},
	}
	for _, input := range inputs {
		res, err := In(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestOverlap(t *testing.T) {
	inputs := []res{
		{[]interface{}{[]interface{}{1, 2, 3}, []interface{}{2, 3, 4}}, true, false},
		{[]interface{}{[]interface{}{true}, []interface{}{true, false}}, true, false},
		{[]interface{}{[]interface{}{1, 2.0, 3.000000}, []interface{}{2, 3, 4}}, true, false},
		{[]interface{}{[]interface{}{"1", "2", "3"}, []interface{}{"2", "3", "4"}}, true, false},
		{[]interface{}{[]interface{}{1, 2, 3}, []interface{}{6, 5, 4}}, false, false},
		{[]interface{}{[]interface{}{true}, []interface{}{false}}, false, false},
		{[]interface{}{[]interface{}{1, 2.0, 3.000000}, []interface{}{6, 5, 4}}, false, false},
		{[]interface{}{[]interface{}{"1", "2", "3"}, []interface{}{"6", "5", "4"}}, false, false},
		{[]interface{}{[]interface{}{1, 2, 3}, []interface{}{2, 3, 4}, []interface{}{2, 3, 4}}, false, true},
		{[]interface{}{[]interface{}{1, 2, 3}, 2}, false, true},
		{[]interface{}{[]interface{}{1, 2, []interface{}{4}}, []interface{}{3, 4, 5}}, false, true},
	}
	for _, input := range inputs {
		res, err := Overlap(input.params...)
		if input.err {
			if err == nil {
				t.Errorf("input: %v, shoud have errors but got none", input.params)
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestBetween(t *testing.T) {
	inputs := []res{
		{[]interface{}{100.0, 10.0, 1000.0}, true, false},
		{[]interface{}{100.0, 10, 1000.0}, true, false},
		{[]interface{}{"b", "a", "c"}, true, false},
		{[]interface{}{time.Now(), time.Now().Add(-100 * time.Second), time.Now()}, true, false},

		{[]interface{}{time.Now(), time.Now(), time.Now()}, false, false},
		{[]interface{}{100.0, 101.0, 1000.0}, false, false},

		{[]interface{}{100.0, "10", 1000.0}, false, true},
		{[]interface{}{"100", "10", 1000.0}, false, true},
		{[]interface{}{"b", "c"}, false, true},
		{[]interface{}{"b", "c", "b", "c"}, false, true},
		{[]interface{}{100.0, 10.0, []int{1000}}, false, true},
	}
	for _, input := range inputs {
		res, err := Between(input.params...)
		if input.err {
			if err == nil {
				t.Errorf("input: %v, shoud have errors but got none", input.params)
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestAnd(t *testing.T) {
	inputs := []res{
		{[]interface{}{true, true}, true, false},
		{[]interface{}{true, true, true, true}, true, false},

		{[]interface{}{true, false, true, true}, false, false},
		{[]interface{}{true, false}, false, false},
		{[]interface{}{false, false}, false, false},

		{[]interface{}{1, false}, false, true},
	}
	fn := AndOr{ModeAnd}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestOr(t *testing.T) {
	inputs := []res{
		{[]interface{}{true, false}, true, false},
		{[]interface{}{true, true}, true, false},
		{[]interface{}{true, true, true, true}, true, false},
		{[]interface{}{true, false, true, true}, true, false},

		{[]interface{}{false, false}, false, false},

		{[]interface{}{1, false}, false, true},
	}
	fn := AndOr{ModeOr}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestNot(t *testing.T) {
	inputs := []res{
		{[]interface{}{false}, true, false},

		{[]interface{}{true}, false, false},

		{[]interface{}{true, false, true, true}, false, true},
		{[]interface{}{1, false}, false, true},
		{[]interface{}{1}, false, true},
	}
	for _, input := range inputs {
		res, err := Not(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestEqual(t *testing.T) {
	inputs := []res{
		{[]interface{}{100.0, 100.0}, true, false},
		{[]interface{}{100, 100.0}, true, false},
		{[]interface{}{now, now}, true, false},
		{[]interface{}{[]interface{}{now}, []interface{}{now}}, true, false},
		{[]interface{}{true, true}, true, false},
		{[]interface{}{"100", "100", "100", "100"}, true, false},

		{[]interface{}{100, 100, 200}, false, false},

		{[]interface{}{"200"}, false, true},
	}
	for _, input := range inputs {
		res, err := Equal{}.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Errorf("input: %v, shoud have errors but got none", input.params)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("input: %v, shoud not have error but got %s", input.params, err.Error())
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestNotEqual(t *testing.T) {
	inputs := []res{
		{[]interface{}{100.0, 200.0}, true, false},
		{[]interface{}{true, false}, true, false},
		{[]interface{}{[]interface{}{true}, []interface{}{false}}, true, false},
		{[]interface{}{"100", "200", "300", "400"}, true, false},

		{[]interface{}{"100", "100", "100", "100"}, false, false},
		{[]interface{}{100, 100.0}, false, false},
		{[]interface{}{100.0, 200.0, 200}, false, false},

		{[]interface{}{100, 100, "100"}, true, true},
		{[]interface{}{"200"}, true, true},
		{[]interface{}{100, 100.1, "100"}, true, true},
		{[]interface{}{[]interface{}{100.0}, 200.0, 200}, true, true},
		{[]interface{}{[]interface{}{true}, []interface{}{false}, false}, true, true},
	}
	for _, input := range inputs {
		res, err := NotEqual(input.params...)
		if input.err {
			if err == nil {
				t.Errorf("input: %v, shoud have errors but got none", input.params)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("input: %v, shoud not have error but got %s", input.params, err.Error())
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestTypeVersion(t *testing.T) {
	inputs := []res{
		{[]interface{}{"1.1.1"}, true, false},
		{[]interface{}{[]interface{}{"2.7.1", "2.8.0"}}, true, false},
		{[]interface{}{[]interface{}{"2.7.1", "2.8.0.1.1.1.1.1.1.1"}}, true, false},
		{[]interface{}{[]interface{}{"2.7.1", "9999.8.0.1.1.1.1.1.1.9999"}}, true, false},

		{[]interface{}{[]interface{}{"2.7.1", "9999.8.0.1.1.1.1.1.1.1.9999"}}, true, true},
		{[]interface{}{[]interface{}{"2.7.1", "2.8.0.1.1.1.1.1.1.10000"}}, true, true},
		{[]interface{}{[]interface{}{"2.7.1", "2.8.0.1.1.1.1.abc.1.1"}}, true, true},
	}
	fn := TypeVersion{}
	for _, input := range inputs {
		_, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
	}
}

func TestTime(t *testing.T) {
	now := time.Now().UTC()
	format1 := "2006-01-02 15:04:05"
	format2 := "2006-01-02"
	now1, _ := time.Parse(format1, now.Format(format1))
	now2, _ := time.Parse(format2, now.Format(format2))
	inputs := []res{
		{[]interface{}{format1, now.Format(format1)}, now1, false},
		{[]interface{}{format2, now.Format(format2)}, now2, false},
		{[]interface{}{format1, []interface{}{now.Format(format1), now.Format(format1)}}, []time.Time{now1, now1}, false},
		{[]interface{}{"2006", []interface{}{now.Format(format1), now.Format(format1)}}, nil, true},
		{[]interface{}{"2016", []interface{}{now.Format(format1), now.Format(format1)}}, nil, true},
		{[]interface{}{"2006", 2016}, nil, true},
	}
	fn := TypeTime{}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if list, ok := res.([]time.Time); ok {
			for i := 0; i < len(list); i++ {
				if in := input.result.([]time.Time)[i]; in != list[i] {
					t.Errorf("input: %v wanna: %v, got: %v", input.params, in, list[i])
				}

			}
		} else {
			if input.result != res {
				t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
			}
		}
	}
}

func TestDefaultTime(t *testing.T) {
	now := time.Now().UTC()
	format1 := "2006-01-02 15:04:05"
	format2 := "2006-01-02"
	now1, _ := time.Parse(format1, now.Format(format1))
	inputs := []res{
		{[]interface{}{now.Format(format1)}, now1, false},
		{[]interface{}{[]interface{}{now.Format(format1), now.Format(format1)}}, []time.Time{now1, now1}, false},
		{[]interface{}{[]interface{}{now.Format(format1), []interface{}{now.Format(format1)}}}, nil, true},
		{[]interface{}{[]interface{}{[]interface{}{now.Format(format1)}, []interface{}{now.Format(format1)}}}, nil, true},
		{[]interface{}{now.Format(format2)}, nil, true},
		{[]interface{}{"2006", []interface{}{now.Format(format1), now.Format(format1)}}, nil, true},
	}
	fn := TypeTime{DefaultTimeFormat}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if list, ok := res.([]time.Time); ok {
			for i := 0; i < len(list); i++ {
				if in := input.result.([]time.Time)[i]; in != list[i] {
					t.Errorf("input: %v wanna: %v, got: %v", input.params, in, list[i])
				}

			}
		} else {
			if input.result != res {
				t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
			}
		}
	}
}

func TestDefaultDate(t *testing.T) {
	now := time.Now().UTC()
	format1 := "2006-01-02"
	format2 := "2006-01-02 15:04:05"
	now1, _ := time.Parse(format1, now.Format(format1))
	inputs := []res{
		{[]interface{}{now.Format(format1)}, now1, false},
		{[]interface{}{[]interface{}{now.Format(format1), now.Format(format1)}}, []time.Time{now1, now1}, false},
		{[]interface{}{[]interface{}{now.Format(format1), []interface{}{now.Format(format1)}}}, nil, true},
		{[]interface{}{[]interface{}{[]interface{}{now.Format(format1)}, []interface{}{now.Format(format1)}}}, nil, true},
		{[]interface{}{now.Format(format2)}, nil, true},
		{[]interface{}{"2006", []interface{}{now.Format(format1), now.Format(format1)}}, nil, true},
	}
	fn := TypeTime{DefaultDateFormat}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if list, ok := res.([]time.Time); ok {
			for i := 0; i < len(list); i++ {
				if in := input.result.([]time.Time)[i]; in != list[i] {
					t.Errorf("input: %v wanna: %v, got: %v", input.params, in, list[i])
				}

			}
		} else {
			if input.result != res {
				t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
			}
		}
	}
}

func TestAdd(t *testing.T) {
	inputs := []res{
		{[]interface{}{100, 100}, 200.0, false},
		{[]interface{}{100, 100.0}, 200.0, false},
		{[]interface{}{100, 100.0, 200.0}, 400.0, false},

		{[]interface{}{100}, .0, true},
		{[]interface{}{100, "100"}, .0, true},
	}
	fn := SuccessiveBinaryOperator{ModeAdd}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestDivide(t *testing.T) {
	inputs := []res{
		{[]interface{}{100, 100}, 1.0, false},

		{[]interface{}{100, 0}, .0, true},
		{[]interface{}{100}, .0, true},
		{[]interface{}{100, "100"}, .0, true},
	}
	fn := BinaryOperator{ModeDivide}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	inputs := []res{
		{[]interface{}{100, 0.0}, true, false},
		{[]interface{}{100, int32(10)}, true, false},
		{[]interface{}{100, 200}, false, false},

		{[]interface{}{100, "0.0"}, false, true},
		{[]interface{}{100, []interface{}{100, 200, 300}}, false, true},
	}
	fn := Compare{ModeGreaterThan}
	for _, input := range inputs {
		res, err := fn.Eval(input.params...)
		if input.err {
			if err == nil {
				t.Error("shoud have errors but got none")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
		}
		if input.result != res {
			t.Errorf("input: %v wanna: %v, got: %v", input.params, input.result, res)
		}
	}
}
