package acidtab

import (
	"fmt"
	"strconv"
)

// PrintAsFloat prints n as a float with with at most 1 digit after the decimal
// point.
func PrintAsFloat(perc int) PrintAsFunc {
	return func(v interface{}) string {
		var f float64
		switch vv := v.(type) {
		default:
			// TODO: maybe add error return to PrintAsFunc?
			panic(fmt.Sprintf("acidtab.PrintAsFloat: not a float but %T: %[1]v", v))

		case float64:
			f = vv
		case float32:
			f = float64(vv)
		}

		if f < 1 {
			return fmt.Sprintf("%.1f", f)[1:]
		}
		return fmt.Sprintf("%0.0f", v)
	}
}

// PrintAsNum prints n as a number with , as thousands separators.
func PrintAsNum(n interface{}) string {
	// TODO: allow configuring this.
	// There's also "indian style" where the grouping is different, but full
	// locale parsing isn't really a goal here.
	sep := ','
	var s string

	switch nn := n.(type) {
	default:
		panic(fmt.Sprintf("acidtab.PrintAsNum: unsupported type: %T: %[1]v", nn))

	// Not really numbers, but just allow it.
	case string:
		s = nn
	case []byte:
		s = string(nn)

	case int:
		s = strconv.FormatInt(int64(nn), 10)
	case int8:
		s = strconv.FormatInt(int64(nn), 10)
	case int16:
		s = strconv.FormatInt(int64(nn), 10)
	case int32:
		s = strconv.FormatInt(int64(nn), 10)
	case int64:
		s = strconv.FormatInt(nn, 10)
	case uint:
		s = strconv.FormatUint(uint64(nn), 10)
	case uint8:
		s = strconv.FormatUint(uint64(nn), 10)
	case uint16:
		s = strconv.FormatUint(uint64(nn), 10)
	case uint32:
		s = strconv.FormatUint(uint64(nn), 10)
	case uint64:
		s = strconv.FormatUint(nn, 10)
	case float32:
		s = strconv.FormatFloat(float64(nn), 'f', 0, 32)
	case float64:
		s = strconv.FormatFloat(nn, 'f', 0, 64)
	}

	if len(s) < 4 {
		return s
	}

	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	var out []rune
	for i := range b {
		if i > 0 && i%3 == 0 && sep > 1 {
			out = append(out, sep)
		}
		out = append(out, rune(b[i]))
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}
