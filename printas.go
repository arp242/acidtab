package acidtab

import (
	"fmt"
	"strconv"
)

// FormatAsFloat prints n as a float with the given percision.
//
// Use perc=0 to round to the nearest natural number.
func FormatAsFloat(perc int) FormatAsFunc {
	return func(v interface{}) string {
		var f float64
		switch vv := v.(type) {
		default:
			// TODO: maybe add error return to FormatAsFunc?
			panic(fmt.Sprintf("acidtab.FormatAsFloat: not a float but %T: %[1]v", v))

		case float64:
			f = vv
		case float32:
			f = float64(vv)
		}

		if f < 1 {
			return fmt.Sprintf("%."+strconv.Itoa(perc)+"f", f)[1:]
		}
		return fmt.Sprintf("%0."+strconv.Itoa(perc)+"f", v)
	}
}

// FormatAsNum prints n as a number with , as thousands separators.
func FormatAsNum() FormatAsFunc {
	return func(v interface{}) string {
		// TODO: allow configuring this.
		// There's also "indian style" where the grouping is different, but full
		// locale parsing isn't really a goal here.
		sep := ','
		var s string

		switch vv := v.(type) {
		default:
			panic(fmt.Sprintf("acidtab.FormatAsNum: unsupported type: %T: %[1]v", vv))

		// Not really numbers, but just allow it.
		case string:
			s = vv
		case []byte:
			s = string(vv)

		case int:
			s = strconv.FormatInt(int64(vv), 10)
		case int8:
			s = strconv.FormatInt(int64(vv), 10)
		case int16:
			s = strconv.FormatInt(int64(vv), 10)
		case int32:
			s = strconv.FormatInt(int64(vv), 10)
		case int64:
			s = strconv.FormatInt(vv, 10)
		case uint:
			s = strconv.FormatUint(uint64(vv), 10)
		case uint8:
			s = strconv.FormatUint(uint64(vv), 10)
		case uint16:
			s = strconv.FormatUint(uint64(vv), 10)
		case uint32:
			s = strconv.FormatUint(uint64(vv), 10)
		case uint64:
			s = strconv.FormatUint(vv, 10)
		case float32:
			s = strconv.FormatFloat(float64(vv), 'f', 0, 32)
		case float64:
			s = strconv.FormatFloat(vv, 'f', 0, 64)
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
}
