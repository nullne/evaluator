package evaluator

import (
	queue "container/list"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/nullne/evaluator/function"
)

var (
	ErrUnexpectedEnd        = errors.New("unexpected end")
	ErrNilInput             = errors.New("nil input")
	ErrLeftOverText         = errors.New("left over text")
	ErrUnmatchedParenthesis = errors.New("unmatched parenthesis")
)

// string without quoted is variable name
type varString string

func (q varString) String() string {
	return string(q)
}

// dynamic types for i are string, qString, int, float64, list
type sexp struct {
	// type of i must NOT be sexp
	i interface{}
}

func (exp sexp) evaluate(pf ParamsFunc) (interface{}, error) {
	if l, isList := exp.i.(list); isList {
		if len(l) == 0 {
			return l, nil
		}
		isFunc := false
		if name, ok := l[0].i.(varString); ok {
			if _, err := function.Get(string(name)); err == nil {
				isFunc = true
			}
		}

		params := make([]interface{}, 0, len(l))
		tl := l
		if isFunc {
			tl = l[1:]
		}
		for _, p := range tl {
			v, err := p.evaluate(pf)
			if err != nil {
				return nil, err
			}
			params = append(params, v)
		}
		if isFunc {
			fn, _ := function.Get(string(l[0].i.(varString)))
			return fn.Eval(params...)
		} else {
			return append(make([]interface{}, 0, len(params)), params...), nil
		}
	} else {
		if val, ok := exp.i.(varString); ok {
			return pf(string(val))
		}
		return exp.i, nil
	}
}

func parse(exp string) (sexp, error) {
	data := []byte(exp)
	tokens := queue.New()
ss:
	for i := 0; i < len(data); {
		advance, token, err := scan(data[i:])
		if err != nil {
			return sexp{}, err
		}
		i += advance
		if t, ok := token.(byte); ok && t == ')' {
			ins := queue.New()
			for e := tokens.Back(); e != nil; e = tokens.Back() {
				tokens.Remove(e)
				if v, ok := e.Value.(byte); ok && v == '(' {
					exps := make(list, 0, ins.Len())
					for e := ins.Back(); e != nil; e = e.Prev() {
						if p, ok := e.Value.(sexp); ok {
							exps = append(exps, p)
						} else {
							exps = append(exps, sexp{e.Value})
						}
					}
					tokens.PushBack(sexp{exps})
					continue ss
				}
				ins.PushBack(e.Value)
			}
			return sexp{}, ErrUnmatchedParenthesis
		}
		tokens.PushBack(token)
	}

	if tokens.Len() == 0 {
		return sexp{}, ErrNilInput
	} else if tokens.Len() != 1 {
		return sexp{}, ErrLeftOverText
	}
	if exp, ok := tokens.Back().Value.(sexp); ok {
		if l, ok := exp.i.(list); ok && len(l) == 0 {
			return sexp{}, ErrNilInput
		}
		return exp, nil
	}
	return sexp{tokens.Back().Value}, nil
}

func (s sexp) String() string {
	return fmt.Sprintf("%v", s.i)
}

func (s sexp) dump(i int) {
	fmt.Printf("%*s%v: ", i*3, "", reflect.TypeOf(s.i))
	if l, isList := s.i.(list); isList {
		fmt.Println(len(l), "elements")
		for _, e := range l {
			e.dump(i + 1)
		}
	} else {
		fmt.Println(s.i)
	}
}

type list []sexp

func (l list) String() string {
	if len(l) == 0 {
		return "[]"
	}
	b := fmt.Sprintf("[%v", l[0])
	for _, s := range l[1:] {
		b = fmt.Sprintf("%s %v", b, s)
	}
	return b + "]"
}

func scan(data []byte) (advance int, token interface{}, err error) {
	length := len(data)
	start := 0
	for width := 0; start < length; start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}
	if start >= length {
		return start, nil, nil
	}

	if b := data[start]; b == '\'' || b == '"' || b == '`' {
		advance, token, err = scanStringWithQuotesStriped(data[start:])
		return start + advance, string(token.([]byte)), err
	}

	for width, i := 0, start; i < length; i += width {
		if b := data[i]; b == ')' || b == '(' {
			if i == start {
				return start + 1, data[i], nil
			} else {
				return i, convert(data[start:i]), nil
			}
		}
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) {
			return i, convert(data[start:i]), nil
		}
	}
	return len(data), convert(data[start:]), nil
}

func convert(data []byte) interface{} {
	if v, err := strconv.Atoi(string(data)); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(string(data), 64); err == nil {
		return v
	}
	return varString(data)
}

// scanStringWithQuotesStriped scan string surrounded with ', " or something like this, a single character
func scanStringWithQuotesStriped(data []byte) (advance int, token []byte, err error) {
	length := len(data)
	if length == 0 {
		return 0, nil, nil
	}
	// escape character positions
	ecps := make([]int, 0, len(data))
	defer func() {
		// remove escape character
		if err != nil {
			return
		}
		tokenLen := len(token)
		tmp := make([]byte, 0, tokenLen)
		s := 0
		for _, p := range ecps {
			tmp = append(tmp, token[s:p]...)
			s = p + 1
		}
		if s < tokenLen {
			tmp = append(tmp, token[s:]...)
		}
		token = tmp
	}()

	delim := data[0]
	ecp := -1
	for i := 1; i < length; i++ {
		if data[i] == '\\' {
			ecp = i
		}
		if data[i] == delim {
			escapeEscaped := continuousCharacterCountFromBack(data[:i], '\\')%2 == 0
			if ecp == i-1 && !escapeEscaped {
				ecps = append(ecps, ecp-1)
			} else {
				advance, token = i+1, data[1:i]
				return
			}
		}
	}
	return 0, nil, ErrUnexpectedEnd
}

func continuousCharacterCountFromBack(data []byte, key byte) int {
	length := len(data)
	for i := length; i > 0; i-- {
		if data[i-1] != key {
			return length - i
		}
	}
	return len(data)
}
