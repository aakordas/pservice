package ptlist

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/aakordas/pservice/time_utils"
)

// Invalid is the unsuccessful request response type.
type Invalid struct {
	Status      string `json:"status"`
	Description string `json:"desc"`
}

func init() {
	http.HandleFunc("/ptlist", interval)
}

func validResponse(w http.ResponseWriter, result time_utils.Interval) {
	w.Header().Set(`Content-Type`, `application/json`)
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(result)
	if err != nil {
		http.Error(w, "Fatal error", http.StatusInternalServerError)
		return
	}
}

func errorResponse(w http.ResponseWriter, desc string, status int) {
	w.Header().Set(`Content-Type`, `application/json`)
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(Invalid{
		Status:      "error",
		Description: desc,
	})
	if err != nil {
		http.Error(w, "Fatal error", http.StatusInternalServerError)
		return
	}
}

func errorBadRequest(w http.ResponseWriter, desc string) {
	errorResponse(w, desc, http.StatusBadRequest)
}

// interval is the handler for the /ptlist endpoint. It expects a period query,
// an optional timezone query and two time values. It responds with an array of
// strings, which contains all the times with delta period the period given in
// the query, as long as the two time values are valid (t1 < t2).
func interval(w http.ResponseWriter, r *http.Request) {
	queries, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		errorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	period := queries.Get("period")
	if period == "" {
		errorBadRequest(w, "Missing period")
		return
	}

	every, err := time.ParseDuration(period)

	// I will not try to manufacture a regexp for the clock time, that
	// covers the entire time.ParseDuration specification. If it is a wrong
	// time, it will show up and throw an error when it will get parsed.
	// I will also not try to make something as sleek for the calendar
	// periods. I am a bit lazy, short on time and it would be an overkill
	// :)
	validDatePeriod := time_utils.DayRegexp.MatchString(period) ||
		time_utils.MonthRegexp.MatchString(period) ||
		time_utils.YearRegexp.MatchString(period)
	if !validDatePeriod && err != nil {
		errorBadRequest(w, "Unsupported period")
		return
	}

	loc := queries.Get("tz")
	if loc == "" {
		// "" is also an accepted Location value, but I believe local
		// time is more convenient than UTC.
		loc = "Local"
	}

	location, err := time.LoadLocation(loc)
	if err != nil {
		errorBadRequest(w, "Bad location provided.")
		return
	}

	t1q := queries.Get("t1")
	if t1q == "" {
		errorBadRequest(w, "Missing t1")
		return
	}

	t2q := queries.Get("t2")
	if t2q == "" {
		errorBadRequest(w, "Missing t2")
		return
	}

	t1, err := time_utils.NewTime(location, t1q)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	t2, err := time_utils.NewTime(location, t2q)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	if t1.After(t2) {
		errorBadRequest(w, "t1 after t2")
		return
	}

	var intervals time_utils.Interval
	if validDatePeriod {
		intervals, err = time_utils.CalculateDateIntervals(period, t1, t2, location)
		if err != nil {
			errorBadRequest(w, "Invalid period")
			return
		}
	} else {
		intervals = time_utils.CalculateTimeIntervals(every, t1, t2)
	}

	validResponse(w, intervals)
}
