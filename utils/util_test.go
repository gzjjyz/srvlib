package utils

import (
	"testing"
)

func TestCalcMillionRate(t *testing.T) {
	type args struct {
		base uint32
		rate uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "base=10, rate=1000",
			args: args{
				base: 10,
				rate: 1000,
			},
			want: 11,
		},
		{
			name: "base=1, rate=1000",
			args: args{
				base: 1,
				rate: 1000,
			},
			want: 1,
		},
		{
			name: "base=0, rate=1000",
			args: args{
				base: 0,
				rate: 1000,
			},
			want: 0,
		},
		{
			name: "base=0, rate=1000",
			args: args{
				base: 0,
				rate: 1000,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcMillionRate(tt.args.base, tt.args.rate); got != tt.want {
				t.Errorf("CalcMillionRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcMillionRate64(t *testing.T) {
	type args struct {
		base int64
		rate int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "base=10, rate=1000",
			args: args{
				base: 10,
				rate: 1000,
			},
			want: 11,
		},
		{
			name: "base=1, rate=1000",
			args: args{
				base: 1,
				rate: 1000,
			},
			want: 1,
		},
		{
			name: "base=0, rate=1000",
			args: args{
				base: 0,
				rate: 1000,
			},
			want: 0,
		},
		{
			name: "base=0, rate=1000",
			args: args{
				base: 0,
				rate: 1000,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcMillionRate64(tt.args.base, tt.args.rate); got != tt.want {
				t.Errorf("CalcMillionRate64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcMillionRateBoth64(t *testing.T) {
	type args struct {
		base int64
		rate int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "base=10, rate=1000",
			args: args{
				base: 10,
				rate: 1000,
			},
			want: 11,
		},
		{
			name: "base=1, rate=1000",
			args: args{
				base: 1,
				rate: 1000,
			},
			want: 1,
		},
		{
			name: "base=0, rate=1000",
			args: args{
				base: 0,
				rate: 1000,
			},
			want: 0,
		},
		{
			name: "base=0, rate=1000",
			args: args{
				base: 0,
				rate: 1000,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcMillionRateBoth64(tt.args.base, tt.args.rate); got != tt.want {
				t.Errorf("CalcMillionRateBoth64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcMillionRateRevert(t *testing.T) {
	type args struct {
		base int64
		rate int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "base=10, rate=1000",
			args: args{
				base: 10,
				rate: 1000,
			},
			want: 9,
		},
		{
			name: "base=1, rate=1000",
			args: args{
				base: 1,
				rate: 1000,
			},
			want: 0,
		},
		{
			name: "base=10000, rate=1000",
			args: args{
				base: 10000,
				rate: 1000,
			},
			want: 9000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcMillionRateRevert(tt.args.base, tt.args.rate); got != tt.want {
				t.Errorf("CalcMillionRateRevert() = %v, want %v", got, tt.want)
			}
		})
	}
}
