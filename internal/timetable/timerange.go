package timetable

import (
	"fmt"
	"time"
)

func max(t, u time.Time) time.Time {
	if t.After(u) {
		return t
	}
	return u
}

func min(t, u time.Time) time.Time {
	if t.Before(u) {
		return t
	}
	return u
}

type TimeRange struct {
	Start, End time.Time
}

func newUnboundTimeRange(start time.Time) (TimeRange, error) {
	return newTimeRange(start, time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 0, 0, start.Location()))
}

func newTimeRange(start, end time.Time) (TimeRange, error) {
	if start.IsZero() {
		return TimeRange{}, fmt.Errorf("start is zero time")
	} else if end.IsZero() {
		return TimeRange{}, fmt.Errorf("end is zero time")
	} else if start.After(end) {
		return TimeRange{}, fmt.Errorf("start %v is after end %v", start, end)
	}

	return TimeRange{start, end}, nil
}

func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

func (tr TimeRange) sub(u TimeRange) []TimeRange {
	res := []TimeRange{}
	if u.Start.Compare(tr.Start) == 1 {
		res = append(res, TimeRange{tr.Start, min(u.Start, tr.End)})
	}
	if u.End.Compare(tr.End) == -1 {
		res = append(res, TimeRange{max(u.End, tr.Start), tr.End})
	}
	return res
}

func (tr TimeRange) subMulti(us []TimeRange) []TimeRange {
	if len(us) == 0 {
		return []TimeRange{tr}
	}

	sub := tr.sub(us[0])

	if len(us) > 1 {
		us = us[1:]
	} else {
		us = []TimeRange{}
	}

	if len(sub) == 1 {
		return sub[0].subMulti(us)
	} else if len(sub) == 2 {
		return append([]TimeRange{sub[0]}, sub[1].subMulti(us)...)
	} else {
		return []TimeRange{}
	}
}
