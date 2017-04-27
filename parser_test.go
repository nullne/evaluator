package condition

import (
	"fmt"
	"testing"
)

type res struct {
	data    string
	advance int
	token   string
	err     error
}

func TestFscan(t *testing.T) {
	s := []byte(`
and
	(or
		(= gender male)
		(> (age birthdate) 18)
	)
	(in version "2.1.6" "2.1.7" "2.1.8")
	`)
	fmt.Println(string(s))
	for i := 0; i < len(s); {
		advance, token, _ := scan(s[i:])
		fmt.Println(string(token))
		i += advance
	}
}

func TestFscanStringWithQuotesStriped(t *testing.T) {
	data := []res{
		{``, 0, ``, nil},
		{`strings`, 7, `tring`, nil},
		{`st'rings`, 8, `t'ring`, nil},
		{`"string"`, 8, `string`, nil},
		{`"str'ing"`, 9, `str'ing`, nil},
		{`"str ing"`, 9, `str ing`, nil},
		{`'string'`, 8, `string`, nil},
		{`'str\'ing'`, 10, `str'ing`, nil},
		{`'str\\'ing'`, 7, `str\\`, nil},
		{`'str\\\'ing'`, 12, `str\\'ing`, nil},
		{`'`, 0, ``, ErrUnexpectedEnd},
		{`string`, 0, ``, ErrUnexpectedEnd},
		{`'str\\\'`, 0, ``, ErrUnexpectedEnd},
		{`st'ring`, 0, ``, ErrUnexpectedEnd},
	}
	for _, d := range data {
		advance, token, err := scanStringWithQuotesStriped([]byte(d.data))
		if d.advance != advance ||
			d.token != string(token) ||
			d.err != err {
			t.Errorf("%s, advanced, token, error expected to be  (%d, %s, %v),  but (%d, %s, %v)", d.data, d.advance, d.token, d.err, advance, token, err)
		}
	}
}
