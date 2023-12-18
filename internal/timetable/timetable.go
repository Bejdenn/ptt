package timetable

import (
	"errors"
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
	ID         int
	Start, End time.Time
	Pause      time.Duration
}

func (s Session) Duration() time.Duration {
	return s.End.Sub(s.Start)
}

func (s Session) String() string {
	return fmt.Sprintf("%d: %v - %v", s.ID, s.Start.Format(time.Kitchen), s.End.Format(time.Kitchen))
}

type Timetable struct {
	Sessions []Session
}

func (t *Timetable) String() string {
	s := ""
	for _, v := range t.Sessions {
		s += fmt.Sprintf("%v, ", v)
	}
	return s
}

type ErrStartAfterEnd struct {
	start time.Time
	end   time.Time
}

func (e ErrStartAfterEnd) Error() string {
	return fmt.Sprintf("start (%v) is after end (%v)", e.start.Format(time.UnixDate), e.end.Format(time.UnixDate))
}

func GenerateTimetable(start, end time.Time, pause time.Duration, sessions SessionInfo) (*Timetable, error) {
	tt := Timetable{}

	if start.IsZero() {
		return nil, fmt.Errorf("start is zero time")
	} else if !end.IsZero() && start.After(end) {
		return nil, ErrStartAfterEnd{start, end}
	}

	if sessions.Duration == time.Duration(0) && end.IsZero() {
		// neither duration nor end time given, so default to 6 hours of work
		sessions.Duration = 6 * time.Hour
	}

	if sessions.Duration != time.Duration(0) {
		endTemp := start.Add(sessions.Duration)
		for i := 0; i < int(math.Ceil(sessions.Duration.Minutes()/sessions.SessionLength.Minutes())-1); i++ {
			endTemp = endTemp.Add(pause)
		}

		if end.IsZero() || end.After(endTemp) {
			end = endTemp
		}
	}

	sessionStart := start
	var sessionEnd time.Time

	for i := 0; ; i++ {
		sessionEnd = sessionStart.Add(sessions.SessionLength)

		if !end.IsZero() && sessionEnd.After(end) {
			// optimize rest time to end by filling up the session with a unit that is a fraction of the given unit length
			opt := int(math.Max(0, end.Sub(sessionStart).Minutes())/float64(step)) * step
			if opt <= threshold {
				break
			}

			sessionEnd = sessionStart.Add(time.Minute * time.Duration(opt))
		}

		s := Session{i + 1, sessionStart, sessionEnd, pause}
		tt.Sessions = append(tt.Sessions, s)

		sessionStart = sessionEnd.Add(pause)
	}

	if len(tt.Sessions) == 0 {
		return nil, errors.New("no sessions generated")
	}

	// remove last pause to not exceed the cumulative time, as you probably are not going to do a pause after the last session
	tt.Sessions[len(tt.Sessions)-1].Pause = 0

	return &tt, nil
}
