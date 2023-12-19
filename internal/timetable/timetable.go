package timetable

import (
	"fmt"
	"math"
	"time"
)

var (
	step      = 5
	threshold = step
)

type SessionInfo struct {
	Duration      time.Duration
	SessionLength time.Duration
}

type Session struct {
	ID        int
	TimeRange TimeRange
	Pause     time.Duration
}

func (s Session) Duration() time.Duration {
	return s.TimeRange.End.Sub(s.TimeRange.Start)
}

func (s Session) String() string {
	return fmt.Sprintf("%d: %v - %v", s.ID, s.TimeRange.Start.Format(time.Kitchen), s.TimeRange.End.Format(time.Kitchen))
}

type TimeRange struct {
	Start, End time.Time
}

func NewTimeRange(start, end time.Time, pause, duration, sessionLength time.Duration) (TimeRange, error) {
	if end.IsZero() {
		if duration == time.Duration(0) {
			// neither end nor duration given, so default to 6 hours of work
			return newTimeRangeByDuration(start, pause, 6*time.Hour, sessionLength)
		} else {
			return newTimeRangeByDuration(start, pause, duration, sessionLength)
		}
	} else {
		if duration == time.Duration(0) {
			return newTimeRangeByEnd(start, end, pause, sessionLength)
		} else {
			c1, err := newTimeRangeByDuration(start, pause, duration, sessionLength)
			if err != nil {
				return TimeRange{}, err
			}
			c2, err := newTimeRangeByEnd(start, end, pause, sessionLength)
			if err != nil {
				return TimeRange{}, err
			}
			if c1.End.Before(c2.End) {
				return c1, nil
			} else {
				return c2, nil
			}
		}
	}
}

func newTimeRangeByDuration(start time.Time, pause, duration, sessionLength time.Duration) (TimeRange, error) {
	if start.IsZero() {
		return TimeRange{}, fmt.Errorf("start is zero time")
	} else if sessionLength == time.Duration(0) {
		return TimeRange{}, fmt.Errorf("session length is zero")
	}

	pauses := int(math.Ceil(duration.Minutes()/sessionLength.Minutes()) - 1)
	end := start.Add(duration).Add(time.Minute * time.Duration(int(pause.Minutes())*pauses))
	return newTimeRangeByEnd(start, end, pause, sessionLength)
}

func newTimeRangeByEnd(start, end time.Time, pause, sessionLength time.Duration) (TimeRange, error) {
	if start.IsZero() {
		return TimeRange{}, fmt.Errorf("start is zero time")
	} else if end.IsZero() {
		return TimeRange{}, fmt.Errorf("end is zero time")
	} else if start.After(end) {
		return TimeRange{}, fmt.Errorf("start %v is after end %v", start, end)
	} else if sessionLength == time.Duration(0) {
		return TimeRange{}, fmt.Errorf("session length is zero")
	}

	return TimeRange{start, end}, nil
}

func Generate(tr TimeRange, pause, sessionLength time.Duration) ([]Session, error) {
	sessions := []Session{}

	sessionStart := tr.Start
	var sessionEnd time.Time

	for i := 0; ; i++ {
		sessionEnd = sessionStart.Add(sessionLength)

		if sessionEnd.After(tr.End) {
			// optimize rest time to end by filling up the session with a unit that is a fraction of the given unit length
			opt := int(math.Max(0, tr.End.Sub(sessionStart).Minutes())/float64(step)) * step
			if opt <= threshold {
				break
			}

			sessionEnd = sessionStart.Add(time.Minute * time.Duration(opt))
		}

		s := Session{i + 1, TimeRange{sessionStart, sessionEnd}, pause}
		sessions = append(sessions, s)

		sessionStart = sessionEnd.Add(pause)
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions generated")
	}

	// remove last pause to not exceed the cumulative time, as you probably are not going to do a pause after the last session
	sessions[len(sessions)-1].Pause = 0

	return sessions, nil
}
