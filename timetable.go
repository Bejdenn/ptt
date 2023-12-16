package main

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

type sessionInfo struct {
	duration      time.Duration
	sessionLength time.Duration
}

type session struct {
	id         int
	start, end time.Time
	pause      time.Duration
}

func (s session) String() string {
	return fmt.Sprintf("%d: %v - %v", s.id, s.start.Format(time.Kitchen), s.end.Format(time.Kitchen))
}

type timetable struct {
	sessions []session
}

func (t *timetable) String() string {
	s := ""
	for _, v := range t.sessions {
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

func generateTimetable(start, end time.Time, pause time.Duration, sessions sessionInfo) (*timetable, error) {
	tt := timetable{}

	if start.IsZero() {
		return nil, fmt.Errorf("start is zero time")
	} else if !end.IsZero() && start.After(end) {
		return nil, ErrStartAfterEnd{start, end}
	}

	if sessions.duration == time.Duration(0) && end.IsZero() {
		// neither duration nor end time given, so default to 6 hours of work
		sessions.duration = 6 * time.Hour
	}

	if sessions.duration != time.Duration(0) {
		endTemp := start.Add(sessions.duration)
		for i := 0; i < int(math.Ceil(sessions.duration.Minutes()/sessions.sessionLength.Minutes())-1); i++ {
			endTemp = endTemp.Add(pause)
		}

		if end.IsZero() || end.After(endTemp) {
			end = endTemp
		}
	}

	sessionStart := start
	var sessionEnd time.Time

	for i := 0; ; i++ {
		sessionEnd = sessionStart.Add(sessions.sessionLength)

		if !end.IsZero() && sessionEnd.After(end) {
			// optimize rest time to end by filling up the session with a unit that is a fraction of the given unit length
			opt := int(math.Max(0, end.Sub(sessionStart).Minutes())/float64(step)) * step
			if opt <= threshold {
				break
			}

			sessionEnd = sessionStart.Add(time.Minute * time.Duration(opt))
		}

		s := session{i + 1, sessionStart, sessionEnd, pause}
		tt.sessions = append(tt.sessions, s)

		sessionStart = sessionEnd.Add(pause)
	}

	if len(tt.sessions) == 0 {
		return nil, errors.New("no sessions generated")
	}

	// remove last pause to not exceed the cumulative time, as you probably are not going to do a pause after the last session
	tt.sessions[len(tt.sessions)-1].pause = 0

	return &tt, nil
}
