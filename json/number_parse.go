package json

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func isJSONNumberNonDelim(c byte) bool {
	return c == '-' || c == '+' || c == '.' || c == '_' ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9')
}

func scanJSONNumber(input string) (int, bool) {
	s := input
	n := 0
	if len(s) == 0 {
		return 0, false
	}
	if s[0] == '-' {
		s = s[1:]
		n++
		if len(s) == 0 {
			return 0, false
		}
	}
	switch {
	case s[0] == '0':
		s = s[1:]
		n++
	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		n++
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
	default:
		return 0, false
	}
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		n += 2
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
	}
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		n++
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			n++
			if len(s) == 0 {
				return 0, false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
	}
	if n < len(input) && isJSONNumberNonDelim(input[n]) {
		return 0, false
	}
	return n, true
}

func parseJSONInt(s string, bitSize int) (int64, error) {
	if len(s) != len(strings.TrimSpace(s)) {
		return 0, fmt.Errorf("invalid syntax")
	}
	n, ok := scanJSONNumber(s)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return intFromJSONNumber(s[:n], bitSize)
}

func parseJSONUint(s string, bitSize int) (uint64, error) {
	if len(s) != len(strings.TrimSpace(s)) {
		return 0, fmt.Errorf("invalid syntax")
	}
	n, ok := scanJSONNumber(s)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return uintFromJSONNumber(s[:n], bitSize)
}

func parseJSONFloat(s string, bitSize int) (float64, error) {
	switch s {
	case "NaN":
		return math.NaN(), nil
	case "Infinity":
		return math.Inf(1), nil
	case "-Infinity":
		return math.Inf(-1), nil
	}
	if len(s) != len(strings.TrimSpace(s)) {
		return 0, fmt.Errorf("invalid syntax")
	}
	n, ok := scanJSONNumber(s)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return strconv.ParseFloat(s[:n], bitSize)
}

// numberParts holds the pieces of a valid JSON number.
//
// Integer conversion works from these parts rather than from a float64.
// Rounding a decimal through a float first loses digits above 2^53 and turns a
// value outside the target range into that range's limit, so both a precise
// large value and an overflow would be reported as an ordinary success.
type numberParts struct {
	neg  bool
	intp string
	frac string
	exp  string
}

// parseNumberParts splits a valid JSON number into its parts. It walks the
// input the same way scanJSONNumber does, and additionally keeps the pieces
// needed to rebuild the number as an integer.
func parseNumberParts(input string) (numberParts, bool) {
	var parts numberParts

	s := input
	if len(s) == 0 {
		return numberParts{}, false
	}

	if s[0] == '-' {
		parts.neg = true
		s = s[1:]
		if len(s) == 0 {
			return numberParts{}, false
		}
	}

	switch {
	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		intp := s
		n := 1
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
		parts.intp = intp[:n]

	default:
		return numberParts{}, false
	}

	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		frac := s[1:]
		n := 1
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
		// Trailing zeroes carry no value, so dropping them lets a number such
		// as 215.0 normalize to an integer.
		parts.frac = strings.TrimRight(frac[:n], "0")
	}

	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		exp := s
		n := 0
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			n++
			if len(s) == 0 {
				return numberParts{}, false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
		parts.exp = exp[:n]
	}

	return parts, true
}

// normalizeToIntString rewrites the parts as a plain decimal integer string,
// resolving any exponent. It reports false when the number is not an integer
// or when the exponent moves the value past what an integer can hold.
func normalizeToIntString(n numberParts) (string, bool) {
	intpSize := len(n.intp)
	fracSize := len(n.frac)

	if intpSize == 0 && fracSize == 0 {
		return "0", true
	}

	var exp int
	if len(n.exp) > 0 {
		i, err := strconv.ParseInt(n.exp, 10, 32)
		if err != nil {
			return "", false
		}
		exp = int(i)
	}

	var num []byte
	if exp >= 0 {
		// A positive exponent shifts fraction digits into the integer part,
		// padding with zeroes once the fraction runs out.
		if fracSize > exp {
			return "", false
		}

		// Stop before building a slice that cannot fit an integer anyway.
		const maxDigits = 20 // Max uint64 value has 20 decimal digits.
		if intpSize+exp > maxDigits {
			return "", false
		}

		num = make([]byte, 0, intpSize+exp)
		num = append(num, n.intp...)
		num = append(num, n.frac...)
		for i := 0; i < exp-fracSize; i++ {
			num = append(num, '0')
		}
	} else {
		// A negative exponent shifts digits out of the integer part, so any
		// fraction at all means the value is not an integer.
		if fracSize > 0 {
			return "", false
		}

		index := intpSize + exp
		if index < 0 {
			return "", false
		}

		num = []byte(n.intp)
		// Digits pushed past the decimal point must all be zero, otherwise
		// the value is not an integer.
		for i := index; i < intpSize; i++ {
			if num[i] != '0' {
				return "", false
			}
		}
		num = num[:index]
	}

	if n.neg {
		return "-" + string(num), true
	}
	return string(num), true
}

func intFromJSONNumber(num string, bitSize int) (int64, error) {
	parts, ok := parseNumberParts(num)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	s, ok := normalizeToIntString(parts)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return strconv.ParseInt(s, 10, bitSize)
}

func uintFromJSONNumber(num string, bitSize int) (uint64, error) {
	parts, ok := parseNumberParts(num)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	s, ok := normalizeToIntString(parts)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return strconv.ParseUint(s, 10, bitSize)
}
