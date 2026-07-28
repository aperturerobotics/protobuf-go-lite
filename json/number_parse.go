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

func parseSignedIntLenient(s string, bitSize int) (int64, error) {
	if len(s) != len(strings.TrimSpace(s)) {
		return 0, fmt.Errorf("invalid syntax")
	}
	n, ok := scanJSONNumber(s)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return intFromJSONNumber(s[:n], bitSize)
}

func parseUnsignedIntLenient(s string, bitSize int) (uint64, error) {
	if len(s) != len(strings.TrimSpace(s)) {
		return 0, fmt.Errorf("invalid syntax")
	}
	n, ok := scanJSONNumber(s)
	if !ok {
		return 0, fmt.Errorf("invalid syntax")
	}
	return uintFromJSONNumber(s[:n], bitSize)
}

func parseFloatLenient(s string, bitSize int) (float64, error) {
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

func intFromJSONNumber(num string, bitSize int) (int64, error) {
	if v, err := strconv.ParseInt(num, 10, bitSize); err == nil {
		return v, nil
	}
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) || math.Trunc(f) != f {
		return 0, fmt.Errorf("invalid syntax")
	}
	v := int64(f)
	if bitSize == 32 {
		if v < math.MinInt32 || v > math.MaxInt32 {
			return 0, fmt.Errorf("value out of range")
		}
	}
	return v, nil
}

func uintFromJSONNumber(num string, bitSize int) (uint64, error) {
	if v, err := strconv.ParseUint(num, 10, bitSize); err == nil {
		return v, nil
	}
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) || f < 0 || math.Trunc(f) != f {
		return 0, fmt.Errorf("invalid syntax")
	}
	v := uint64(f)
	if bitSize == 32 && v > math.MaxUint32 {
		return 0, fmt.Errorf("value out of range")
	}
	return v, nil
}
