package function

import (
	"testing"
)

func TestGet(t *testing.T) {
	if err := Regist("foo", foo); err != nil {
		t.Error(err)
	}
	if err := RegistFuncer("foor", foor{}); err != nil {
		t.Error(err)
	}

	if err := Regist("foo", foo); err != ErrFunctionExists {
		t.Error("should exist")
	}
	if err := RegistFuncer("foor", foor{}); err != ErrFunctionExists {
		t.Error("should exist")
	}

	if _, err := Get("foo"); err != nil {
		t.Error(err)
	}
	if _, err := Get("fooooo"); err != ErrNotFound {
		t.Error("should not found")
	}
	_ = Registered()
}

var foo = func(params ...interface{}) (interface{}, error) {
	return nil, nil
}

type foor struct{}

func (f foor) Eval(params ...interface{}) (interface{}, error) {
	return nil, nil
}
