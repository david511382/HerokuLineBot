package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
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
	const (
		DOT   = "."
		COMMA = ","
	)

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
