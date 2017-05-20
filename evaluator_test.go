package condition

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	// stmt := `(in gender ("female" "male"))`
	stmt := `(in 123 (123 456))`
	exp, err := parse(stmt)
	if err != nil {
		t.Error(err)
	}
	vvf := func(name string) (interface{}, error) {
		return "male", nil
	}
	res, err := evaluate(exp, vvf)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
