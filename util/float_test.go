package util

import (
	"testing"
)

func TestFloat_Round(t *testing.T) {
	type args struct {
		v   float64
		exp int32
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"0",
			args{
				v:   1.4,
				exp: 0,
			},
			1,
		},
		{
			"int",
			args{
				v:   12.5,
				exp: 1,
			},
			10,
		},
		{
			"float",
			args{
				v:   5.45,
				exp: -1,
			},
			5.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFloat(tt.args.v)
			if got := f.Round(tt.args.exp).Value(); got != tt.want {
				t.Errorf("FloatRound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFloat_PlusFloat(b *testing.B) {
	dataCount := 1000.0
	arr := []float64{}
	for i := 0.0; i < dataCount; i++ {
		arr = append(arr, 1000-i)
	}
	for i := 0; i < b.N; i++ {
		f := NewFloat(0)
		f.PlusFloat(arr...)
	}
}
