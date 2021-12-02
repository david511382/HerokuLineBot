package util

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type Float decimal.Decimal

func NewFloat(v float64) Float {
	return Float(decimal.NewFromFloat(v))
}

func NewInt64Float(v int64) Float {
	return Float(decimal.NewFromInt(v))
}

func (f Float) GetUtilCompareValue() float64 {
	return f.Value()
}

func (f Float) Value() float64 {
	d := decimal.Decimal(f)
	r, _ := d.Float64()
	return r
}

func (f Float) ToInt() int64 {
	d := decimal.Decimal(f)
	return d.IntPart()
}

// 5.45 ToString(-1) // output: "5.5"
func (f Float) ToString(prec int) string {
	return strconv.FormatFloat(f.Value(), 'f', prec, 64)
}

func (f Float) Minus(vs ...Float) Float {
	d1 := decimal.Decimal(f)
	for _, v := range vs {
		d2 := decimal.Decimal(v)
		d1 = d1.Sub(d2)
	}
	return Float(d1)
}

func (f Float) MinusFloat(vs ...float64) Float {
	for _, v := range vs {
		d := NewFloat(v)
		f = f.Minus(d)
	}
	return f
}

func (f Float) Plus(vs ...Float) Float {
	d1 := decimal.Decimal(f)
	for _, v := range vs {
		d2 := decimal.Decimal(v)
		d1 = d1.Add(d2)
	}
	return Float(d1)
}

func (f Float) PlusFloat(vs ...float64) Float {
	for _, v := range vs {
		d := NewFloat(v)
		f = f.Plus(d)
	}
	return f
}

func (f Float) PlusInt64(vs ...int64) Float {
	for _, v := range vs {
		d := NewInt64Float(v)
		f = f.Plus(d)
	}
	return f
}

func (f Float) Mul(vs ...Float) Float {
	d1 := decimal.Decimal(f)
	for _, v := range vs {
		d2 := decimal.Decimal(v)
		d1 = d1.Mul(d2)
	}
	return Float(d1)
}

func (f Float) MulFloat(vs ...float64) Float {
	for _, v := range vs {
		d := NewFloat(v)
		f = f.Mul(d)
	}
	return f
}

func (f Float) Div(v Float) Float {
	d1 := decimal.Decimal(f)
	d2 := decimal.Decimal(v)
	return Float(d1.Div(d2))
}

func (f Float) DivFloat(v float64) Float {
	return f.Div(NewFloat(v))
}

func (f Float) DivInt64(v int64) Float {
	return f.Div(NewInt64Float(v))
}

// 4捨5入
// 5.45 Round(-1) // output: "5.5"
// 12.5 Round(1) // output: "10"
func (f Float) Round(exp int32) Float {
	d := decimal.Decimal(f)
	d = d.Round(-exp)
	return Float(d)
}

// Floor returns the nearest integer value less than or equal to d.
func (f Float) Floor() Float {
	d := decimal.Decimal(f)
	d = d.Floor()
	return Float(d)
}

// Ceil returns the nearest integer value greater than or equal to d.
func (f Float) Ceil() Float {
	d := decimal.Decimal(f)
	d = d.Ceil()
	return Float(d)
}

func (f Float) Percent(exp int32) Float {
	return f.MulFloat(100)
}
