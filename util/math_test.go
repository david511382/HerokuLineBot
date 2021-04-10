package util

import "testing"

func TestFloatString(t *testing.T) {
	type args struct {
		value         float64
		floatExponent int32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"short",
			args{
				value:         0.01,
				floatExponent: -1,
			},
			"0.0",
		},
		{
			"zero exponent",
			args{
				value:         123.01,
				floatExponent: 0,
			},
			"123",
		},
		{
			"exponent",
			args{
				value:         100.01,
				floatExponent: -1,
			},
			"100.0",
		},
		{
			"thousand",
			args{
				value:         12002334.01,
				floatExponent: -1,
			},
			"12,002,334.0",
		},
		{
			"negative",
			args{
				value:         -111223.01,
				floatExponent: 0,
			},
			"-111,223",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FloatString(tt.args.value, tt.args.floatExponent); got != tt.want {
				t.Errorf("FloatString() = %v, want %v", got, tt.want)
			}
		})
	}
}
