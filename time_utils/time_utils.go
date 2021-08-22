package time_utils

import (
	"time"
)

// Interval is the successful request response type.
type Interval []string

// timeFormat is desired time format. It is of the form "yyyymmddThhmmssZ"
const timeFormat string = "20060102T150405Z"

// NewTime returns a new Time in the specified location and with the specified
// clock.
func NewTime(loc, clock string) (time.Time, error) {
	location, err := time.LoadLocation(loc)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(timeFormat, clock, location)

	return t, err
}

// PrintTime prints the given time in the "yyyymmddThhmmssZ" format.
func PrintTime(t time.Time) string {
	return t.Format(timeFormat)
}

// TODO: Test if I have 12:00, for example, does Add.Truncate return 12 or 13?

// CalculateTimeIntervals returns the Interval elapsed between t1 and t2 in
// period intervals. period should be a duration string valid for
// time.ParseDuration.
func CalculateTimeIntervals(period time.Duration, t1, t2 time.Time) (result Interval) {
	start := t1.Add(period).Truncate(period)
	end := t2.Truncate(period)

	for t := start; t.Before(end); t = t.Add(period) {
		result = append(result, PrintTime(t))
	}
	result = append(result, PrintTime(end))

	return result
}

// addDay returns t+1d
func addDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

// addMonth returns t+1mo
func addMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, 0)
}

// addYear returns t+1y
func addYear(t time.Time) time.Time {
	return t.AddDate(1, 0, 0)
}

// calculateIntervals returns an Interval with all the elapsed intervals between
// start and end, as calculated by addTime.
func calculateIntervals(start, end time.Time, addTime func(time.Time) time.Time) (result Interval) {
	for t := start; t.Before(end); t = addTime(t) {
		result = append(result, PrintTime(t))
	}

	return result[:len(result)-1]
}

// CalculateDateIntervals returns the Interval elapsed between t1 and t2 in
// period intervals. Valid period units are "d", for days, "mo", for months, and
// "y" for years.
func CalculateDateIntervals(period string, t1, t2 time.Time) (result Interval) {
	if period == "1d" {
		result = calculateIntervals(t1, t2, addDay)
	} else if period == "1mo" {
		result = calculateIntervals(t1, t2, addMonth)
	} else if period == "1y" {
		result = calculateIntervals(t1, t2, addYear)
	}

	return result
}
