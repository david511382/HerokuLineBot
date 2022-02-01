package util

import (
	"testing"
)

func TestDescSearch(t *testing.T) {
	arr := []int{222333, 22222, 111, 110, 99, 38, 9, 6, 3, 1}
	type args struct {
		startIndex int
		lastIndex  int
		compareF   func(index int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"standard",
			args{
				0,
				9,
				func(index int) int {
					value := 38
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			5,
		},
		{
			"last",
			args{
				0,
				9,
				func(index int) int {
					value := 1
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			9,
		},
		{
			"first",
			args{
				0,
				9,
				func(index int) int {
					value := 222333
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescSearch(tt.args.startIndex, tt.args.lastIndex, tt.args.compareF); got != tt.want {
				t.Errorf("DescSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	arr := []int{2, 4, 5, 7, 10, 88, 939, 2944, 9993, 222333}
	type args struct {
		startIndex int
		lastIndex  int
		compareF   func(index int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"standard",
			args{
				0,
				9,
				func(index int) int {
					value := 2944
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			7,
		},
		{
			"last",
			args{
				0,
				9,
				func(index int) int {
					value := 222333
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			9,
		},
		{
			"first",
			args{
				0,
				9,
				func(index int) int {
					value := 2
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Search(tt.args.startIndex, tt.args.lastIndex, tt.args.compareF); got != tt.want {
				t.Errorf("DescSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	value := 499
	arr := []int{}
	for i := 0; i < 1000; i++ {
		arr = append(arr, 1000-i)
	}
	for i := 0; i < b.N; i++ {
		DescSearch(0, 999, func(index int) int {
			if value > arr[index] {
				return 1
			} else if value == arr[index] {
				return 0
			}
			return -1
		})
	}
}

func BenchmarkNormalSearch(b *testing.B) {
	value := 499
	arr := []int{}
	for i := 0; i < 1000; i++ {
		arr = append(arr, 1000-i)
	}
	for i := 0; i < b.N; i++ {
		for _, v := range arr {
			if value == v {
				break
			}
		}
	}
}

func TestInsertIndex(t *testing.T) {
	arr := []int{2, 4, 5, 7, 10, 88, 939, 2944, 9993, 222333}
	type args struct {
		startIndex int
		lastIndex  int
		compareF   func(index int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"standard",
			args{
				0,
				9,
				func(index int) int {
					value := 2945
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			8,
		},
		{
			"last",
			args{
				0,
				9,
				func(index int) int {
					value := 2223331
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			10,
		},
		{
			"first",
			args{
				0,
				9,
				func(index int) int {
					value := 1
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			0,
		},
		{
			"equal",
			args{
				0,
				9,
				func(index int) int {
					value := 939
					if value > arr[index] {
						return 1
					} else if value == arr[index] {
						return 0
					}
					return -1
				},
			},
			7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InsertIndex(tt.args.startIndex, tt.args.lastIndex, tt.args.compareF); got != tt.want {
				t.Errorf("InsertIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
