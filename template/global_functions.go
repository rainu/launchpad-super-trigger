package template

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var globalFuncMap map[string]interface{}

func getGlobalFunctions() map[string]interface{} {
	if globalFuncMap == nil {
		globalFuncMap = map[string]interface{}{
			"add": add,
			"sub": sub,
			"mul": mul,
			"div": div,

			"strings_Count":          strings.Count,
			"strings_Contains":       strings.Contains,
			"strings_ContainsAny":    strings.ContainsAny,
			"strings_ContainsRune":   strings.ContainsRune,
			"strings_LastIndex":      strings.LastIndex,
			"strings_IndexByte":      strings.IndexByte,
			"strings_IndexRune":      strings.IndexRune,
			"strings_IndexAny":       strings.IndexAny,
			"strings_LastIndexAny":   strings.LastIndexAny,
			"strings_LastIndexByte":  strings.LastIndexByte,
			"strings_SplitN":         strings.SplitN,
			"strings_SplitAfterN":    strings.SplitAfterN,
			"strings_Split":          strings.Split,
			"strings_SplitAfter":     strings.SplitAfter,
			"strings_Fields":         strings.Fields,
			"strings_FieldsFunc":     strings.FieldsFunc,
			"strings_Join":           strings.Join,
			"strings_HasPrefix":      strings.HasPrefix,
			"strings_HasSuffix":      strings.HasSuffix,
			"strings_Map":            strings.Map,
			"strings_Repeat":         strings.Repeat,
			"strings_ToUpper":        strings.ToUpper,
			"strings_ToLower":        strings.ToLower,
			"strings_ToTitle":        strings.ToTitle,
			"strings_ToUpperSpecial": strings.ToUpperSpecial,
			"strings_ToLowerSpecial": strings.ToLowerSpecial,
			"strings_ToTitleSpecial": strings.ToTitleSpecial,
			"strings_ToValidUTF8":    strings.ToValidUTF8,
			"strings_Title":          strings.Title,
			"strings_TrimLeftFunc":   strings.TrimLeftFunc,
			"strings_TrimRightFunc":  strings.TrimRightFunc,
			"strings_TrimFunc":       strings.TrimFunc,
			"strings_IndexFunc":      strings.IndexFunc,
			"strings_LastIndexFunc":  strings.LastIndexFunc,
			"strings_Trim":           strings.Trim,
			"strings_TrimLeft":       strings.TrimLeft,
			"strings_TrimRight":      strings.TrimRight,
			"strings_TrimSpace":      strings.TrimSpace,
			"strings_TrimPrefix":     strings.TrimPrefix,
			"strings_TrimSuffix":     strings.TrimSuffix,
			"strings_Replace":        strings.Replace,
			"strings_ReplaceAll":     strings.ReplaceAll,
			"strings_EqualFold":      strings.EqualFold,
			"strings_Index":          strings.Index,
			"strings_NewReader":      strings.NewReader,
			"strings_Compare":        strings.Compare,
			"strings_NewReplacer":    strings.NewReplacer,

			"math_Abs":             math.Abs,
			"math_J0":              math.J0,
			"math_Y0":              math.Y0,
			"math_J1":              math.J1,
			"math_Y1":              math.Y1,
			"math_Jn":              math.Jn,
			"math_Yn":              math.Yn,
			"math_Dim":             math.Dim,
			"math_Max":             math.Max,
			"math_Min":             math.Min,
			"math_Erf":             math.Erf,
			"math_Erfc":            math.Erfc,
			"math_Exp":             math.Exp,
			"math_Exp2":            math.Exp2,
			"math_FMA":             math.FMA,
			"math_Log":             math.Log,
			"math_Mod":             math.Mod,
			"math_Pow":             math.Pow,
			"math_Cos":             math.Cos,
			"math_Sin":             math.Sin,
			"math_Tan":             math.Tan,
			"math_Asin":            math.Asin,
			"math_Acos":            math.Acos,
			"math_Atan":            math.Atan,
			"math_Inf":             math.Inf,
			"math_NaN":             math.NaN,
			"math_IsNaN":           math.IsNaN,
			"math_IsInf":           math.IsInf,
			"math_Cbrt":            math.Cbrt,
			"math_Logb":            math.Logb,
			"math_Ilogb":           math.Ilogb,
			"math_Sinh":            math.Sinh,
			"math_Cosh":            math.Cosh,
			"math_Sqrt":            math.Sqrt,
			"math_Tanh":            math.Tanh,
			"math_Acosh":           math.Acosh,
			"math_Asinh":           math.Asinh,
			"math_Atan2":           math.Atan2,
			"math_Atanh":           math.Atanh,
			"math_Expm1":           math.Expm1,
			"math_Floor":           math.Floor,
			"math_Ceil":            math.Ceil,
			"math_Trunc":           math.Trunc,
			"math_Round":           math.Round,
			"math_RoundToEven":     math.RoundToEven,
			"math_Gamma":           math.Gamma,
			"math_Hypot":           math.Hypot,
			"math_Ldexp":           math.Ldexp,
			"math_Log10":           math.Log10,
			"math_Log2":            math.Log2,
			"math_Log1p":           math.Log1p,
			"math_Pow10":           math.Pow10,
			"math_Erfinv":          math.Erfinv,
			"math_Erfcinv":         math.Erfcinv,
			"math_Float32bits":     math.Float32bits,
			"math_Float32frombits": math.Float32frombits,
			"math_Float64bits":     math.Float64bits,
			"math_Float64frombits": math.Float64frombits,
			"math_Signbit":         math.Signbit,
			"math_Copysign":        math.Copysign,
			"math_Nextafter32":     math.Nextafter32,
			"math_Nextafter":       math.Nextafter,
			"math_Remainder":       math.Remainder,

			"strconv_ParseBool":                strconv.ParseBool,
			"strconv_FormatBool":               strconv.FormatBool,
			"strconv_AppendBool":               strconv.AppendBool,
			"strconv_ParseFloat":               strconv.ParseFloat,
			"strconv_ParseUint":                strconv.ParseUint,
			"strconv_ParseInt":                 strconv.ParseInt,
			"strconv_Atoi":                     strconv.Atoi,
			"strconv_FormatFloat":              strconv.FormatFloat,
			"strconv_AppendFloat":              strconv.AppendFloat,
			"strconv_FormatUint":               strconv.FormatUint,
			"strconv_FormatInt":                strconv.FormatInt,
			"strconv_Itoa":                     strconv.Itoa,
			"strconv_AppendInt":                strconv.AppendInt,
			"strconv_AppendUint":               strconv.AppendUint,
			"strconv_Quote":                    strconv.Quote,
			"strconv_AppendQuote":              strconv.AppendQuote,
			"strconv_QuoteToASCII":             strconv.QuoteToASCII,
			"strconv_AppendQuoteToASCII":       strconv.AppendQuoteToASCII,
			"strconv_QuoteToGraphic":           strconv.QuoteToGraphic,
			"strconv_AppendQuoteToGraphic":     strconv.AppendQuoteToGraphic,
			"strconv_QuoteRune":                strconv.QuoteRune,
			"strconv_AppendQuoteRune":          strconv.AppendQuoteRune,
			"strconv_QuoteRuneToASCII":         strconv.QuoteRuneToASCII,
			"strconv_AppendQuoteRuneToASCII":   strconv.AppendQuoteRuneToASCII,
			"strconv_QuoteRuneToGraphic":       strconv.QuoteRuneToGraphic,
			"strconv_AppendQuoteRuneToGraphic": strconv.AppendQuoteRuneToGraphic,
			"strconv_CanBackquote":             strconv.CanBackquote,
			"strconv_Unquote":                  strconv.Unquote,
			"strconv_IsPrint":                  strconv.IsPrint,
			"strconv_IsGraphic":                strconv.IsGraphic,

			"rand_NewSource":   rand.NewSource,
			"rand_New":         rand.New,
			"rand_Int63":       rand.Int63,
			"rand_Uint32":      rand.Uint32,
			"rand_Uint64":      rand.Uint64,
			"rand_Int31":       rand.Int31,
			"rand_Int":         rand.Int,
			"rand_Int63n":      rand.Int63n,
			"rand_Int31n":      rand.Int31n,
			"rand_Intn":        rand.Intn,
			"rand_Float64":     rand.Float64,
			"rand_Float32":     rand.Float32,
			"rand_Perm":        rand.Perm,
			"rand_Read":        rand.Read,
			"rand_NormFloat64": rand.NormFloat64,
			"rand_ExpFloat64":  rand.ExpFloat64,

			"time_Parse":                  time.Parse,
			"time_ParseInLocation":        time.ParseInLocation,
			"time_ParseDuration":          time.ParseDuration,
			"time_Since":                  time.Since,
			"time_Until":                  time.Until,
			"time_Now":                    time.Now,
			"time_Unix":                   time.Unix,
			"time_Date":                   time.Date,
			"time_FixedZone":              time.FixedZone,
			"time_LoadLocation":           time.LoadLocation,
			"time_LoadLocationFromTZData": time.LoadLocationFromTZData,

			"fmt_Sprintf": fmt.Sprintf,
		}
	}

	return globalFuncMap
}

func add(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() + int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() + bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() + float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() + float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() + bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("unknown type for %q (%T)", av, a)
	}
}

func sub(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() - int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() - bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() - float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() - float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() - bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("unknown type for %q (%T)", av, a)
	}
}

func mul(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() * int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) * bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() * bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) * bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() * float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() * float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() * bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("unknown type for %q (%T)", av, a)
	}
}

func div(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() / int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) / bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() / bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) / bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() / float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() / float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() / bv.Float(), nil
		default:
			return nil, fmt.Errorf("unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("unknown type for %q (%T)", av, a)
	}
}
