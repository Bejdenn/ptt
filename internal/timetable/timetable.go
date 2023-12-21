package timetable

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"text/tabwriter"
	"time"
)

const (
	TimeOnlyNoSeconds = "15:04"
	maxDuration       = 1<<63 - 1 // separate constant because the standard library does not export it
)

type Session struct {
	ID        int
	TimeRange TimeRange
	Pause     time.Duration
}

func Generate(start, end time.Time, pause, duration, sessionLength time.Duration, excludes []TimeRange) (SessionSlice, error) {
	var (
		sessions []Session
		err      error
		utr, tr  TimeRange
	)

	utr, err = newUnboundTimeRange(start)
	if err != nil {
		return nil, err
	}

	if end.IsZero() {
		if duration == time.Duration(0) {
			// neither end nor duration given, so default to 6 hours of work
			sessions, err = generate(utr, pause, sessionLength, excludes, 6*time.Hour)
		} else {
			sessions, err = generate(utr, pause, sessionLength, excludes, duration)
		}
	} else {
		tr, err = newTimeRange(start, end)
		if err != nil {
			return nil, err
		}

		if duration == time.Duration(0) {
			sessions, err = generate(tr, pause, sessionLength, excludes)
		} else {
			c1, err := generate(utr, pause, sessionLength, excludes, duration)
			if err != nil {
				return nil, err
			}
			c2, err := generate(tr, pause, sessionLength, excludes)
			if err != nil {
				return nil, err
			}

			if len(c1) == 0 {
				sessions = c2
			} else if len(c2) == 0 {
				sessions = c1
			} else {
				if c1[len(c1)-1].TimeRange.End.Before(c2[len(c2)-1].TimeRange.End) {
					sessions = c1
				} else {
					sessions = c2
				}
			}
		}
	}

	sort.Sort(SessionSlice(sessions))

	if len(sessions) > 0 {
		// remove last pause to not exceed the cumulative time, as you probably are not going to do a pause after the last session
		sessions[len(sessions)-1].Pause = 0
	}

	for i := 0; i < len(sessions); i++ {
		sessions[i].ID = i + 1
	}

	return sessions, err
}

type SessionSlice []Session

func (s SessionSlice) Len() int {
	return len(s)
}

func (s SessionSlice) Less(i, j int) bool {
	return s[i].TimeRange.Start.Before(s[j].TimeRange.Start)
}

func (s SessionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SessionSlice) String() string {
	buf := bytes.NewBuffer([]byte{})

	const padding = 3
	w := tabwriter.NewWriter(buf, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ID\tStart\tEnd\tDuration\tPause\tCumulated Work\tCumulated Time")

	cumulatedWork := time.Duration(0)
	cumulatedTime := time.Duration(0)
	for _, u := range s {
		cumulatedWork += u.TimeRange.Duration()
		cumulatedTime += u.TimeRange.Duration() + u.Pause
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t\n", u.ID, u.TimeRange.Start.Format(TimeOnlyNoSeconds), u.TimeRange.End.Format(TimeOnlyNoSeconds), u.TimeRange.Duration(), u.Pause, cumulatedWork, cumulatedTime)
	}

	w.Flush()

	return buf.String()
}

func generate(tr TimeRange, pause, sessionLength time.Duration, excludes []TimeRange, duration ...time.Duration) ([]Session, error) {
	var d time.Duration
	if len(duration) < 1 {
		d = maxDuration
	} else {
		d = duration[0]
	}

	sessions := []Session{}
	for _, slot := range tr.subMulti(excludes) {
		t := TimeRange{slot.Start, time.Time{}}

		for d > 0 && t.Start.Before(slot.End) {
			t.End = t.Start.Add(time.Minute * time.Duration(math.Min(math.Min(sessionLength.Minutes(), d.Minutes()), slot.End.Sub(t.Start).Minutes())))

			d -= t.Duration()

			sessions = append(sessions, Session{0, t, pause})

			t.Start = t.End.Add(pause)
		}
	}

	return sessions, nil
}
