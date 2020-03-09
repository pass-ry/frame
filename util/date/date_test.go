package date

import (
	"testing"
	"time"
)

func TestIsHolidayByTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "元旦",
			args: args{
				t: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Now().Location()),
			},
			want: true,
		},
		{
			name: "工作日",
			args: args{
				t: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Now().Location()),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHolidayByTime(tt.args.t); got != tt.want {
				t.Errorf("IsHolidayByTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
