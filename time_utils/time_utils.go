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
func NewTime(location *time.Location, clock string) (time.Time, error) {
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
	start := t1.Round(period)
	end := t2.Round(period)

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

func addMonthWithLocation(t time.Time, location *time.Location) func(time.Time) time.Time {
	loc := location
	return func(t time.Time) time.Time {
		return addMonth(t, loc)
	}
}

// roundHour returns t.Round(time.Hour).Hour()
func roundHour(t time.Time) int {
	return t.Round(time.Hour).Hour()
}

// addMonth returns t+1mo, always on the last day of the month
func addMonth(t time.Time, location *time.Location) time.Time {
	firstDay := time.Date(t.Year(), t.Month()+1, 1, roundHour(t), 0, 0, 0, location)
	return firstDay.AddDate(0, 1, -1)
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

	return result
}

// CalculateDateIntervals returns the Interval elapsed between t1 and t2 in
// period intervals. Valid period units are "d", for days, "mo", for months, and
// "y" for years.
func CalculateDateIntervals(period string, t1, t2 time.Time, location *time.Location) (result Interval) {
	// Rounding "inspiration" from https://stackoverflow.com/a/36830387
	if period == "1d" {
		start := t1.Round(time.Hour)
		end := t2.Round(time.Hour)
		result = calculateIntervals(start, end, addDay)
	} else if period == "1mo" {
		firstDay := time.Date(t1.Year(), t1.Month(), 1, roundHour(t1), 0, 0, 0, location)
		start := firstDay.AddDate(0, 1, -1)
		addMo := addMonthWithLocation(start, location)
		result = calculateIntervals(start, t2, addMo)
	} else if period == "1y" {
		// The last day is always 31/12...
		start := time.Date(t1.Year(), 12, 31, roundHour(t1), 0, 0, 0, location)
		result = calculateIntervals(start, t2, addYear)
	}

	return result
}
