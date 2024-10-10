package bankholidays

import (
	"testing"
	"time"
)

func TestIsHoliday(t *testing.T) {
	holidays, err := FetchHolidays()
	if err != nil {
		t.Fatalf("Failed to load holidays: %v", err)
	}
	bh := BankHolidays{
		Holidays: holidays,
		Weekend:  map[time.Weekday]bool{time.Saturday: true, time.Sunday: true},
	}
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !bh.IsHoliday(date, EnglandAndWales) {
		t.Errorf("Expected January 1, 2024 to be a holiday")
	}
}
