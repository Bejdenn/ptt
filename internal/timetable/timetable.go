package timetable

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/Bejdenn/timerange"
)

const (
	maxDuration = 1<<63 - 1 // separate constant because the standard library does not export it
)

type SessionLength struct {
	Min time.Duration
	Max time.Duration
}

func NewSessionLength(min, max time.Duration) (SessionLength, error) {
	if min > max {
		return SessionLength{}, fmt.Errorf("min must not be greater than max in session-length")
	}
	return SessionLength{min, max}, nil
}

type Session struct {
	ID        int
	TimeRange timerange.TimeRange
	Pause     time.Duration
}

func Generate(start, end time.Time, pause, duration time.Duration, sessionLength SessionLength, excludes []timerange.TimeRange) (SessionSlice, error) {
	var (
		sessions []Session
		err      error
	)

	if duration < -1 {
		return nil, fmt.Errorf("duration cannot be negative. only exception is '-1ns' to indicate absence of the duration value")
	}

	utr, err := timerange.NewUnbound(start)
	if err != nil {
		return nil, err
	}

	if end.IsZero() {
		if duration == -1 {
			return nil, fmt.Errorf("neither END or DURATION is set. consider passing either of the values via the respective flag")
		} else {
			sessions, err = generate(utr, pause, sessionLength, excludes, duration)
		}
	} else {
		tr, err := timerange.New(start, end)
		if err != nil {
			return nil, err
		}

		if duration == -1 {
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
				return nil, fmt.Errorf("duration-based generation did not result in any sessions, so end-based generation is ignored")
			} else if len(c2) == 0 {
				return nil, fmt.Errorf("end-time-based generation did not result in any sessions, so duration-based generation is ignored")
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

	for i := range sessions {
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
	fmt.Fprintln(w, "ID\tStart\tEnd\tDuration\tPause\tCumulated Work")

	cumulatedWork := time.Duration(0)
	for _, u := range s {
		cumulatedWork += u.TimeRange.Duration()
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", u.ID, u.TimeRange.Start.Format(timerange.TimeOnlyNoSeconds), u.TimeRange.End.Format(timerange.TimeOnlyNoSeconds), u.TimeRange.Duration(), u.Pause, cumulatedWork)
	}

	w.Flush()

	return buf.String()
}

func generate(tr timerange.TimeRange, pause time.Duration, sessionLength SessionLength, excludes []timerange.TimeRange, duration ...time.Duration) ([]Session, error) {
	var d time.Duration
	if len(duration) < 1 {
		d = maxDuration
	} else {
		d = duration[0]
	}

	sessions := []Session{}
	for _, slot := range tr.SubMulti(excludes) {
		t := timerange.TimeRange{Start: slot.Start, End: time.Time{}}

		for d > 0 && t.Start.Before(slot.End) {
			t.End = t.Start.Add(time.Minute * time.Duration(min(sessionLength.Max.Minutes(), d.Minutes(), slot.End.Sub(t.Start).Minutes())))

			d -= t.Duration()

			sessions = append(sessions, Session{0, t, pause})

			t.Start = t.End.Add(pause)
		}

		if len(sessions) > 0 {
			// remove last pause to not exceed the cumulative time, as you probably are not going to do a pause after the last session
			sessions[len(sessions)-1].Pause = 0
		}
	}

	var out SessionSlice
	for _, s := range sessions {
		if s.TimeRange.Duration() >= sessionLength.Min {
			out = append(out, s)
		}
	}
	sessions = out

	return sessions, nil
}
