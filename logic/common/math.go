package common

import (
	"heroku-line-bot/logic/common/domain"
	"heroku-line-bot/util"
)

func FloatToInt(v float64) int64 {
	return util.FloatToInt(v, domain.FLOAT_EXPONENT)
}

func FloatMinus(v1, v2 float64) float64 {
	return util.FloatMinus(v1, v2, domain.FLOAT_EXPONENT)
}

func FloatPlus(v1, v2 float64) float64 {
	return util.FloatPlus(v1, v2, domain.FLOAT_EXPONENT)
}

func FloatRound(v float64, exp int32) float64 {
	return util.FloatRound(v, exp, domain.FLOAT_EXPONENT)
}

func SafeRate(fraction, denominator float32) float64 {
	return SafeRate64(float64(fraction), float64(denominator))
}

func SafeRate64(fraction, denominator float64) float64 {
	return SafeRate64Exponent(fraction, denominator, domain.FRONT_END_FLOAT_EXPONENT)
}

func SafeRate64Exponent(fraction, denominator float64, f int32) float64 {
	return util.SafeRate64Exponent(fraction, denominator, f, domain.FLOAT_EXPONENT)
}

func SafeDivision64(fraction, denominator float64, f int32) float64 {
	return util.SafeDivision64(fraction, denominator, f, domain.FLOAT_EXPONENT)
}

func PercentAt(value float64, f int32) float64 {
	return util.PercentAt(value, f, domain.FLOAT_EXPONENT)
}

func Percent(value float64) float64 {
	return PercentAt(value, domain.FRONT_END_FLOAT_EXPONENT)
}

func UnlimitSum(a1, r float64) float64 {
	return util.UnlimitSum(a1, r, domain.FLOAT_EXPONENT)
}
