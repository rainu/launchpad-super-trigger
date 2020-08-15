package expressions

import (
	"regexp"
	"strconv"
	"strings"
)

func BuildEqExpressionFn(expr string) func([]byte) bool {
	isNumericValue := false
	fExpValue, err := strconv.ParseFloat(expr, 64)
	if err == nil {
		isNumericValue = true
	}

	return func(value []byte) bool {
		if isNumericValue {
			fValue, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return false
			}

			return fValue == fExpValue
		}

		return expr == string(value)
	}
}

func BuildNeExpressionFn(expr string) func([]byte) bool {
	isNumericValue := false
	fExpValue, err := strconv.ParseFloat(expr, 64)
	if err == nil {
		isNumericValue = true
	}

	return func(value []byte) bool {
		if isNumericValue {
			fValue, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return false
			}

			return fValue != fExpValue
		}

		return expr != string(value)
	}
}
func BuildLtExpressionFn(expr float64) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue < expr
		}

		return false
	}
}

func BuildLteExpressionFn(expr float64) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue <= expr
		}

		return false
	}
}

func BuildGtExpressionFn(expr float64) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue > expr
		}

		return false
	}
}

func BuildGteExpressionFn(expr float64) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue >= expr
		}

		return false
	}
}

func BuildMatchExpressionFn(expr string) func([]byte) bool {
	re := regexp.MustCompile(expr)
	return func(value []byte) bool {
		return re.Match(value)
	}
}

func BuildNotMatchExpressionFn(expr string) func([]byte) bool {
	re := regexp.MustCompile(expr)
	return func(value []byte) bool {
		return !re.Match(value)
	}
}

func BuildContainsExpressionFn(expr string) func([]byte) bool {
	return func(value []byte) bool {
		return strings.Contains(string(value), expr)
	}
}

func BuildNotContainsExpressionFn(expr string) func([]byte) bool {
	return func(value []byte) bool {
		return !strings.Contains(string(value), expr)
	}
}
