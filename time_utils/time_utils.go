package time_utils

import (
	"regexp"
	"strconv"
	"time"
)

// Interval is the successful request response type.
type Interval []string

// timeFormat is desired time format. It is of the form "yyyymmddThhmmssZ"
const timeFormat string = `20060102T150405Z`

// HourRegexp is the regular expression that matches a period of hours.
var HourRegexp *regexp.Regexp = regexp.MustCompile(`[[:digit:]]+h`)

// DayRegexp is the regular expression that matches a period of days.
var DayRegexp *regexp.Regexp = regexp.MustCompile(`[[:digit:]]+d`)

// MonthRegexp is the regular expression that matches a period of months.
var MonthRegexp *regexp.Regexp = regexp.MustCompile(`[[:digit:]]+mo`)

// YearRegexp is the regular expression that matches a period of years.
var YearRegexp *regexp.Regexp = regexp.MustCompile(`[[:digit:]]+y`)

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

// CalculateTimeIntervals returns the Interval elapsed between t1 and t2 in
// period intervals. period should be a duration string of the form "nh", where
// n is the number of hours.
func CalculateTimeIntervals(period time.Duration, t1, t2 time.Time) (result Interval) {
	start := t1.Round(time.Hour)
	end := t2.Round(time.Hour)

	for t := start; t.Before(end); t = t.Add(period) {
		result = append(result, PrintTime(t))
	}

	return result
}

//addDayWithInterval returns a function that adds interval hours to t.
func addDayWithInterval(t time.Time, interval int) func(time.Time) time.Time {
	return func(t time.Time) time.Time {
		return addDay(t, interval)
	}
}

// addDay returns t+1d
func addDay(t time.Time, interval int) time.Time {
	return t.AddDate(0, 0, interval)
}

// addMonthWithLocation returns a function that adds interval months to t.
func addMonthWithLocation(t time.Time, interval int, location *time.Location) func(time.Time) time.Time {
	loc := location
	return func(t time.Time) time.Time {
		return addMonth(t, interval, loc)
	}
}

// roundHour returns t.Round(time.Hour).Hour()
func roundHour(t time.Time) int {
	return t.Round(time.Hour).Hour()
}

// addMonth returns t+1mo, always on the last day of the month
func addMonth(t time.Time, interval int, location *time.Location) time.Time {
	firstDay := time.Date(t.Year(), t.Month()+1, 1, roundHour(t), 0, 0, 0, location)
	return firstDay.AddDate(0, interval, -1)
}

// addYearWithInterval returns a function that adds interval years to t.
func addYearWithInterval(t time.Time, interval int) func(time.Time) time.Time {
	return func(t time.Time) time.Time {
		return addYear(t, interval)
	}
}

// addYear returns t+1y.
func addYear(t time.Time, interval int) time.Time {
	return t.AddDate(interval, 0, 0)
}

// calculateIntervals returns an Interval with all the elapsed intervals between
// start and end, as calculated by addTime.
func calculateIntervals(start, end time.Time, interval int, addTime func(time.Time) time.Time) (result Interval) {
	for t := start; t.Before(end); t = addTime(t) {
		result = append(result, PrintTime(t))
	}

	return result
}

// CalculateDateIntervals returns the Interval elapsed between t1 and t2 in
// period intervals. Valid period units are "d", for days, "mo", for months, and
// "y" for years.
func CalculateDateIntervals(period string, t1, t2 time.Time, location *time.Location) (result Interval, err error) {
	// A bit lazy, but it works...
	var intervalRegexp *regexp.Regexp = regexp.MustCompile(`[^dmy]*`)
	interval, err := strconv.Atoi(intervalRegexp.FindString(period))
	if err != nil {
		return nil, err
	}

	// Rounding "inspiration" from https://stackoverflow.com/a/36830387
	if DayRegexp.MatchString(period) {
		start := t1.Round(time.Hour)
		end := t2.Round(time.Hour)
		addDa := addDayWithInterval(start, interval)
		result = calculateIntervals(start, end, interval, addDa)
	} else if MonthRegexp.MatchString(period) {
		firstDay := time.Date(t1.Year(), t1.Month(), 1, roundHour(t1), 0, 0, 0, location)
		start := firstDay.AddDate(0, 1, -1)
		addMo := addMonthWithLocation(start, interval, location)
		result = calculateIntervals(start, t2, interval, addMo)
	} else if YearRegexp.MatchString(period) {
		// The last day is always 31/12...
		start := time.Date(t1.Year(), 12, 31, roundHour(t1), 0, 0, 0, location)
		addYe := addYearWithInterval(start, interval)
		result = calculateIntervals(start, t2, interval, addYe)
	}

	return result, nil
}
