package service

import (
	"fmt"
	"time"
)


func resolveTimezone(tz string) (*time.Location, error) {
	if tz == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(tz)
}

func buildWindow(date time.Time, startStr, endStr string, loc *time.Location) (time.Time, time.Time, error) {
	var startH, startM, endH, endM int
	if _, err := fmt.Sscanf(startStr, "%d:%d", &startH, &startM); err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_time format")
	}
	if _, err := fmt.Sscanf(endStr, "%d:%d", &endH, &endM); err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_time format")
	}

	windowStart := time.Date(date.Year(), date.Month(), date.Day(), startH, startM, 0, 0, loc)
	windowEnd := time.Date(date.Year(), date.Month(), date.Day(), endH, endM, 0, 0, loc)

	return windowStart, windowEnd, nil
}