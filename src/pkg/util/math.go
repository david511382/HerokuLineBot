package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

const (
	DOT   = "."
	COMMA = ","
)

func UnlimitSum(a1, r float64) float64 {
	dA1 := decimal.NewFromFloat(a1)
	dR := decimal.NewFromFloat(r)
	d1 := decimal.NewFromFloat(1)

	dSum := dA1.Div(d1.Sub(dR))

	result, _ := dSum.Float64()

	return result
}

func FloatString(value float64, floatExponent int32) string {
	format := "%"
	if floatExponent <= 0 {
		format += "0" + DOT + strconv.Itoa(int(-floatExponent))
	}
	format += "f"

	symbol := ""
	if isNegative := value < 0; isNegative {
		value = -value
		symbol = "-"
	}

	raw := fmt.Sprintf(format, value)
	dotIndex := strings.Index(raw, DOT)
	if dotIndex == -1 {
		dotIndex = len(raw)
	}
	shiffStartIndex := dotIndex - 3

	fromIndex := shiffStartIndex % 3
	results := make([]string, 0)

	if fromIndex > 0 {
		s := raw[0:fromIndex]
		results = append(results, s)
	}

	for from := fromIndex; from < shiffStartIndex; from += 3 {
		to := from + 3
		s := raw[from:to]
		results = append(results, s)
	}

	toIndex := len(raw)
	if shiffStartIndex < 0 {
		shiffStartIndex = 0
	}
	s := raw[shiffStartIndex:toIndex]
	results = append(results, s)

	result := strings.Join(results, COMMA)
	result = symbol + result

	return result
}

func StringThousandComma(numStr string) string {
	splitNum := strings.Split(numStr, ".")
	catchInt := splitNum[0]
	if len(catchInt) <= 3 {
		return numStr
	}

	count := 0
	comma := ""

	for i := len(catchInt) - 1; i >= 0; i-- {
		count++
		comma = comma + string(catchInt[i])
		if count%3 == 0 {
			comma = comma + ","
		}
	}

	turnBack := ""
	for i := len(comma) - 1; i >= 0; i-- {
		turnBack = turnBack + string(comma[i])
	}
	splitNum[0] = strings.Trim(turnBack, ",")

	newStr := strings.Join(splitNum, ".")
	return newStr
}

// 去除 i 後 cut 位數
func IntCut[T constraints.Integer](i T, cut int) T {
	return i / T(math.Pow10(cut))
}
