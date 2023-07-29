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
}

type timetable struct {
	sessions            []session
	totalWork, totalDur time.Duration
}

type ErrStartAfterEnd struct {
	start time.Time
	end   time.Time
}

func (e ErrStartAfterEnd) Error() string {
	return fmt.Sprintf("start (%v) is after end (%v)", e.start.Format(time.TimeOnly), e.end.Format(time.TimeOnly))
}

func generateTimetable(start, end time.Time, pausePattern []time.Duration, sessions sessionInfo) (*timetable, error) {
	tt := timetable{}

	if start.IsZero() {
		return nil, fmt.Errorf("start is zero time")
	} else if start.After(end) {
		return nil, ErrStartAfterEnd{start, end}
	}

	if sessions.duration == time.Duration(0) && end.IsZero() {
		// neither duration nor end time given, so default to 6 hours of work
		sessions.duration = 6 * time.Hour
	}

	if sessions.duration != time.Duration(0) {
		endTemp := start.Add(sessions.duration)
		for i := 0; i < int(math.Ceil(sessions.duration.Minutes()/sessions.sessionLength.Minutes())-1); i++ {
			endTemp = endTemp.Add(pausePattern[i%len(pausePattern)])
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

		tt.totalWork += sessionEnd.Sub(sessionStart)
		tt.sessions = append(tt.sessions, session{i + 1, sessionStart, sessionEnd})

		sessionStart = sessionEnd.Add(pausePattern[i%len(pausePattern)])
	}

	if len(tt.sessions) == 0 {
		return nil, errors.New("no sessions generated")
	}

	tt.totalDur = tt.sessions[len(tt.sessions)-1].end.Sub(start)

	return &tt, nil
}
