package time_utils

import (
	"testing"
	"time"
)

func TestCalculateTimeIntervals(t *testing.T) {
	type testCase struct {
		Name string	// The name of the test
		Period time.Duration
		T1, T2 time.Time
	}

	durationsStrings := []string{"1h", "2h", "30h"}
	durations := make([]time.Duration, len(durationsStrings))
	for i := 0; i < len(durationsStrings); i = i+1 {
		var err error
		durations[i], err = time.ParseDuration(durationsStrings[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	location, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}

	type test struct{
		Got testCase
		Want Interval
	}

	got := []test{
		test{
			Got: testCase{
				Name: "1 hour",
				Period: durations[0],
				T1: time.Date(2020, 1, 1, 1, 0, 3, 0, location),
				T2: time.Date(2020, 1, 1, 5, 5, 3, 0, location),
			},
			Want: Interval{"20200101T010000Z", "20200101T020000Z",
				"20200101T030000Z", "20200101T040000Z",
				"20200101T050000Z",
			},
		},
		test{
			Got: testCase{
				Name: "2 hours",
				Period: durations[1],
				T1: time.Date(2020, 1, 1, 1, 2, 0, 0, location),
				T2: time.Date(2020, 1, 1, 5, 2, 5, 0, location),
			},
			Want: Interval{"20200101T010000Z", "20200101T030000Z",
				"20200101T050000Z",
			},
		},
		test{
			Got: testCase{
				Name: "More hours than t1-t2",
				Period: durations[2],
				T1: time.Date(2020, 1, 1, 1, 2, 0, 0, location),
				T2: time.Date(2020, 1, 1, 5, 2, 5, 0, location),
			},
			Want: Interval{"20200101T010000Z"},
		},
		test{
			Got: testCase{
				Name: "Equal times",
				Period: durations[0],
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{},
		},
	}

	for _, te := range got {
		g := CalculateTimeIntervals(te.Got.Period, te.Got.T1, te.Got.T2)

		if !g.Equal(&te.Want) {
			t.Errorf("%s: got %v, want %v\n", te.Got.Name, g, te.Want);
		}
	}
}

func TestCalculateDateIntervals(t *testing.T) {
	type testCase struct {
		Name string	// The name of the test
		Period string
		T1, T2 time.Time
	}

	durations := []string{"1d", "1mo", "1y", "3d", "3mo", "3y", "30y",}

	location, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		Got testCase
		Want Interval
	}

	got := []test{
		test{
			Got: testCase{
				Name: "1 day",
				Period: "1d",
				T1: time.Date(2020, 1, 1, 1, 2, 3, 0, location),
				T2: time.Date(2020, 1, 5, 1, 2, 3, 0, location),
			},
			Want: Interval{"20200101T010000Z", "20200102T010000Z",
				"20200103T010000Z", "20200104T010000Z"},
		},
		test{
			Got: testCase{
				Name: "1 month",
				Period: "1mo",
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2020, 5, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{"20200131T030000Z", "20200229T030000Z",
				"20200331T030000Z", "20200430T030000Z"},
		},
		test{
			Got: testCase{
				Name: "1 year",
				Period: "1y",
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2025, 1, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{"20201231T030000Z", "20211231T030000Z",
				"20221231T030000Z", "20231231T030000Z",
				"20241231T030000Z"},
		},
		test{
			Got: testCase{
				Name: "3 days",
				Period: "3d",
				T1: time.Date(2020, 1, 1, 2, 3, 4, 0, location),
				T2: time.Date(2020, 1, 5, 2, 3, 4, 0, location),
			},
			Want: Interval{"20200101T020000Z", "20200104T020000Z"},
		},
		test{
			Got: testCase{
				Name: "3 months",
				Period: "3mo",
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2020, 5, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{"20200131T030000Z", "20200430T030000Z"},
		},
		test{
			Got: testCase{
				Name: "3 years",
				Period: "3y",
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2025, 1, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{"20201231T030000Z", "20231231T030000Z"},
		},
		test{
			Got: testCase{
				Name: "More years than t1-t2",
				Period: "30y",
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2021, 1, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{"20201231T030000Z"},
		},
		test{
			Got: testCase{
				Name: "t1=t2",
				Period: "1d",
				T1: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
				T2: time.Date(2020, 1, 2, 3, 4, 5, 0, location),
			},
			Want: Interval{},
		},
	}

	for _, te := range got {
		var err error = nil
		g, err := CalculateDateIntervals(te.Got.Period, te.Got.T1, te.Got.T2, location)

		if !g.Equal(&te.Want) || err != nil {
			t.Errorf("%s: got %v, want %v\n", te.Got.Name, g, te.Want);
		}
	}
}
