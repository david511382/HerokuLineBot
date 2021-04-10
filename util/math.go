package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func UnlimitSum(a1, r float64, floatExponent int32) float64 {
	dA1 := decimal.NewFromFloatWithExponent(a1, floatExponent)
	dR := decimal.NewFromFloatWithExponent(r, floatExponent)
	d1 := decimal.NewFromFloatWithExponent(1, floatExponent)

	dSum := dA1.Div(d1.Sub(dR))

	result, _ := dSum.Float64()

	return result
}

func FloatToInt(v float64, floatExponent int32) int64 {
	d := decimal.NewFromFloatWithExponent(v, floatExponent)
	return d.IntPart()
}

func FloatMinus(v1, v2 float64, floatExponent int32) float64 {
	d1 := decimal.NewFromFloatWithExponent(v1, floatExponent)
	d2 := decimal.NewFromFloatWithExponent(v2, floatExponent)
	r := d1.Sub(d2)
	f, _ := r.Float64()
	return f
}

func FloatPlus(v1, v2 float64, floatExponent int32) float64 {
	d1 := decimal.NewFromFloatWithExponent(v1, floatExponent)
	d2 := decimal.NewFromFloatWithExponent(v2, floatExponent)
	r := d1.Add(d2)
	f, _ := r.Float64()
	return f
}

func FloatRound(v float64, exp int32, floatExponent int32) float64 {
	d := decimal.NewFromFloatWithExponent(v, floatExponent)
	d = d.Round(-exp)
	f, _ := d.Float64()
	return f
}

func SafeRate64Exponent(fraction, denominator float64, f int32, floatExponent int32) float64 {
	if denominator == 0 {
		return 0
	}

	fractionV := decimal.NewFromFloatWithExponent(fraction, floatExponent)
	denominatorV := decimal.NewFromFloatWithExponent(denominator, floatExponent)
	r := fractionV.Div(denominatorV)
	v, _ := r.Float64()

	return PercentAt(v, f, floatExponent)
}

func SafeDivision64(fraction, denominator float64, f int32, floatExponent int32) float64 {
	if denominator == 0 {
		return 0
	}

	d1 := decimal.NewFromFloatWithExponent(fraction, floatExponent)
	d2 := decimal.NewFromFloatWithExponent(denominator, floatExponent)
	r := d1.Div(d2)
	r = r.Round(-f)
	v, _ := r.Float64()
	return v
}

func PercentAt(value float64, f int32, floatExponent int32) float64 {
	percent := decimal.NewFromFloatWithExponent(100, floatExponent)
	valueDec := decimal.NewFromFloatWithExponent(value, floatExponent)
	r := valueDec.Mul(percent)
	r = r.Round(-f)
	v, _ := r.Float64()
	return v
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
