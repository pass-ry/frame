package date

import (
	"testing"
	"time"
)

func TestIsForceStopByTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{
				t: time.Now(),
			},
			want: false,
		},
		{
			name: "in control",
			args: args{
				t: time.Date(2020, time.January, 24, 1, 0, 0, 0, time.Now().Location()),
			},
			want: true,
		},
		{
			name: "in control2",
			args: args{
				t: time.Date(2020, time.January, 30, 1, 0, 0, 0, time.Now().Location()),
			},
			want: true,
		},
		{
			name: "out of range",
			args: args{
				t: time.Date(2020, time.January, 31, 1, 0, 0, 0, time.Now().Location()),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsForceStopByTime(tt.args.t); got != tt.want {
				t.Errorf("IsForceStopByTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
