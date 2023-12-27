package time_test

import (
	"testing"
	"time"

	tfmt "github.com/michaelahli/octopus/svcutils/time"
)

func TestTruncateToStartOfDay(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Errorf("TruncateToStartOfDay() err %v", err)
	}

	sod, err := time.ParseInLocation("2006-01-02", time.Now().In(loc).Format("2006-01-02"), loc)
	if err != nil {
		t.Errorf("TruncateToStartOfDay() err %v", err)
	}

	tests := []struct {
		name string
		args time.Time
		want time.Time
	}{
		{
			name: "Success",
			args: time.Now().In(loc),
			want: sod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tfmt.TruncateToStartOfDay(tt.args)
			if !got.Equal(tt.want) {
				t.Errorf("TruncateToStartOfDay() got %v want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateToEndOfDay(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Errorf("TruncateToEndOfDay() err %v", err)
	}

	sod, err := time.ParseInLocation("2006-01-02", time.Now().In(loc).AddDate(0, 0, 1).Format("2006-01-02"), loc)
	if err != nil {
		t.Errorf("TruncateToEndOfDay() err %v", err)
	}

	tests := []struct {
		name string
		args time.Time
		want time.Time
	}{
		{
			name: "Success",
			args: time.Now().In(loc),
			want: sod.Add(-time.Second),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tfmt.TruncateToEndOfDay(tt.args)
			if !got.Equal(tt.want) {
				t.Errorf("TruncateToEndOfDay() got %v want %v", got, tt.want)
			}
		})
	}
}
