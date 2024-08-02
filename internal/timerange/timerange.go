package timerange

import (
	"fmt"
	"strings"
	"time"
)

const TimeOnlyNoSeconds = "15:04"

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

func NewUnbound(start time.Time) (TimeRange, error) {
	return New(start, time.Date(start.Year(), start.Month(), start.Day()+1, 23, 59, 0, 0, start.Location()))
}

func New(start, end time.Time) (TimeRange, error) {
	if start.IsZero() {
		return TimeRange{}, fmt.Errorf("start is zero time")
	} else if end.IsZero() {
		return TimeRange{}, fmt.Errorf("end is zero time")
	} else if start.After(end) {
		return TimeRange{}, fmt.Errorf("start (%v) is after end (%v)", start, end)
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

func (tr TimeRange) SubMulti(us []TimeRange) []TimeRange {
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
		return sub[0].SubMulti(us)
	} else if len(sub) == 2 {
		return append([]TimeRange{sub[0]}, sub[1].SubMulti(us)...)
	} else {
		return []TimeRange{}
	}
}

// Parse parses a list of time ranges from a list of strings.
// The individual strings are expected to be in the format "HH:MM-HH:MM".
// As the time is internally parsed without knowing the date, the date is taken from a reference time, ref.
func Parse(ref time.Time, values []string) ([]TimeRange, error) {
	excludes := make([]TimeRange, 0, len(values))
	for _, exclude := range values {
		parts := strings.Split(exclude, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("exclude did not contain hyphen delimiter: %s", exclude)
		}

		start, err := time.Parse(TimeOnlyNoSeconds, parts[0])
		if err != nil {
			return nil, fmt.Errorf("could not parse start time: %v", err)
		}

		end, err := time.Parse(TimeOnlyNoSeconds, parts[1])
		if err != nil {
			return nil, fmt.Errorf("could not parse end time: %v", err)
		}

		excludes = append(excludes, TimeRange{Start: normalize(ref, start), End: normalize(ref, end)})
	}
	return excludes, nil
}

func normalize(ref, t time.Time) time.Time {
	return time.Date(ref.Year(), ref.Month(), ref.Day(), t.Hour(), t.Minute(), 0, 0, ref.Location())
}
