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

	if period != "1h" && period != "1d" && period != "1mo" && period != "1y" {
		errorBadRequest(w, "Invalid period")
		return
	}

	location := queries.Get("tz")
	if location == "" {
		// "" is also an accepted Location value, but I believe local
		// time is more convenient than UTC.
		location = "Local"
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

	if period == "1h" {
		// I have already checked earlier that period has one of the
		// valid values, so ParseDuration shouldn't return an error.
		every, _ := time.ParseDuration(period)

		validResponse(w, time_utils.CalculateTimeIntervals(every, t1, t2))
	} else {
		validResponse(w, time_utils.CalculateDateIntervals(period, t1, t2))
	}
}
