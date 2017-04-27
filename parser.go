package condition

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	ErrUnexpectedEnd error = errors.New("unexpected end")
)

func Parse(in string) (*unit, error) {
	return nil, nil
}

type unit struct {
	function string
	parent   *unit
	params   []string
}

func (u unit) Eval() string {
	return ""
}

func scan(data []byte) (advance int, token []byte, err error) {
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
		return start + advance, token, err
	}

	defer func() {
		tmp := make([]byte, 0, len(token))
		tmp = append(tmp, token...)
		token = tmp
	}()

	for width, i := 0, start; i < length; i += width {
		if b := data[i]; b == ')' || b == '(' {
			if i == start {
				return start + 1, data[i : i+1], nil
			} else {
				return i, data[start:i], nil
			}
		}
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) {
			return i, data[start:i], nil
		}
	}
	return len(data), data[start:], nil
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
