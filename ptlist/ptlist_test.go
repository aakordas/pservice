// Rather lousy testing because of lack of time.

package ptlist

import(
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// buildQuery builds and returns a query with the provided arguments. The
// returned query is of the form
// ?period=period&tz=timezone&t1=t1&t2=t2. timezone is the only argument that
// can be nil or empty. All the other arguments are mandatory. On error it
// returns the empty string "".
func buildQuery(period, timezone, t1, t2 string) string {
	query := "?"
	periodQuery := "period="
	if period != "" {
		periodQuery += period
		query += periodQuery+"&"
	}

	timezoneQuery := "tz="
	timezoneQuery += timezone
	query += timezoneQuery+"&"

	t1Query := "t1="
	if t1 != "" {
		t1Query += t1
		query += t1Query+"&"
	}

	t2Query := "t2="
	if t2 != "" {
		t2Query += t2
		query += t2Query
	}

	return query
}

func TestResponses(t *testing.T) {
	type testCase struct {
		Name string	// The name of the test
		HTTPMethod string // Just GET for now
		Query string
	}

	type test struct {
		Got testCase
		Want int	// Just the HTTP return code
	}

	got := []test{
		test{
			Got: testCase{
				Name: "1 hour no timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1h", "", "20200101T010000Z", "20200101T050000Z"),
			},
			Want: http.StatusOK,
		},
		test{
			Got: testCase{
				Name: "1 day no timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1d", "", "20200101T010000Z", "20200105T010000Z"),
			},
			Want: http.StatusOK,
		},
		test{
			Got: testCase{
				Name: "3 hours Athens timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("3h", "Europe/Athens", "20200101T010000Z", "20200101T050000Z"),
			},
			Want: http.StatusOK,
		},
		test{
			Got: testCase{
				Name: "1 month Athens timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1mo", "Europe/Athens", "20200101T010000Z", "20200501T010000Z"),
			},
			Want: http.StatusOK,
		},
		test{
			Got: testCase{
				Name: "1 hour Los Angeles timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1h", "America/Los_Angeles", "20200101T010000Z", "20200101T050000Z"),
			},
			Want: http.StatusOK,
		},
		test{
			Got: testCase{
				Name: "1 year Los Angeles timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1y", "America/Los_Angeles", "20200101T010000Z", "20250101T010000Z"),
			},
			Want: http.StatusOK,
		},
		test{
			Got: testCase{
				Name: "3 hours invalid timezone",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("3h", "Africa/Wakanda", "20200101T010000Z", "20200101T050000Z"),
			},
			Want: http.StatusBadRequest,
		},
		test{
			Got: testCase{
				Name: "Missing period",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("", "", "20200101T010000Z", "20200101T010000Z"),
			},
			Want: http.StatusBadRequest,
		},
		test{
			Got: testCase{
				Name: "Missing t1",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1h", "", "", "20200101T010000Z"),
			},
			Want: http.StatusBadRequest,
		},
		test{
			Got: testCase{
				Name: "Missing t2",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1h", "", "20200101T010000Z", ""),
			},
			Want: http.StatusBadRequest,
		},
		test{
			Got: testCase{
				Name: "Bad t1",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1h", "", "20200101T010000ZQ", "20200101T010000Z"),
			},
			Want: http.StatusBadRequest,
		},
		test{
			Got: testCase{
				Name: "Bad t2",
				HTTPMethod: http.MethodGet,
				Query: buildQuery("1h", "", "20200101T010000Z", "20200101T010000ZQ"),
			},
			Want: http.StatusBadRequest,
		},
	}

	uri := "localhost:8282/"
	endpoint := "ptlist"

	for _, g := range got {
		req := httptest.NewRequest(g.Got.HTTPMethod, uri+endpoint+g.Got.Query, nil)
		w := httptest.NewRecorder()
		interval(w, req)

		resp := w.Result()
		// The error can be ignored, because
		// httptest.ResponseRecorder.Result() is guaranteed to return
		// only io.EOF.
		body, _ := io.ReadAll(resp.Body)
		if len(body) == 0 {
			t.Errorf("%s returned an empty body", g.Got.Query)
		}

		if resp.StatusCode != g.Want {
			t.Errorf("%s returned with code %d instead of %d", g.Got.Query, resp.StatusCode, g.Want)
		}
	}
}
