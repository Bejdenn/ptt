package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
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

func generateTimetable(start, end time.Time, pausePattern string, sessions sessionInfo) (*timetable, error) {
	tt := timetable{}

	pauseDurations := []time.Duration{}
	for _, pps := range strings.Split(pausePattern, "-") {
		ppi, err := time.ParseDuration(pps)
		if err != nil {
			return nil, fmt.Errorf("could not parse pause pattern string: %v", err)
		}
		pauseDurations = append(pauseDurations, ppi)
	}

	if sessions.duration == time.Duration(0) && end.IsZero() {
		// neither duration nor end time given, so default to 6 hours of work
		sessions.duration = time.Duration(6 * time.Hour)
	}

	if sessions.duration != time.Duration(0) {
		endTemp := start.Add(sessions.duration)
		for i := 0; i < int(math.Ceil(sessions.duration.Minutes()/sessions.sessionLength.Minutes())-1); i++ {
			endTemp = endTemp.Add(pauseDurations[i%len(pauseDurations)])
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

		sessionStart = sessionEnd.Add(pauseDurations[i%len(pauseDurations)])
	}

	if len(tt.sessions) == 0 {
		return nil, errors.New("no sessions generated")
	}

	tt.totalDur = tt.sessions[len(tt.sessions)-1].end.Sub(start)

	return &tt, nil
}
