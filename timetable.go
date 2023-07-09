package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
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
		sessions.duration = time.Duration(6) * time.Hour
	}

	sessionStart := start
	var sessionEnd time.Time

	for i := 0; ; i++ {
		sessionEnd = sessionStart.Add(sessions.sessionLength)

		optimized := false
		if !end.IsZero() && sessionEnd.After(end) {
			// optimize rest time to end by filling up the session with a unit that is a fraction of the given unit length
			opt := int(math.Max(0, end.Sub(sessionStart).Minutes())/float64(step)) * step
			if opt <= step {
				break
			}

			sessionEnd = sessionStart.Add(time.Minute * time.Duration(opt))
			optimized = true
		}

		if sessions.duration != time.Duration(0) && tt.totalWork+sessionEnd.Sub(sessionStart) > sessions.duration {
			opt := int(math.Max(0, (sessions.duration-tt.totalWork).Minutes())/float64(step)) * step
			if opt <= step {
				break
			}

			sessionEnd = sessionStart.Add(time.Minute * time.Duration(opt))
			optimized = true
		}

		tt.totalWork += sessionEnd.Sub(sessionStart)
		tt.sessions = append(tt.sessions, session{i + 1, sessionStart, sessionEnd})

		// this flag is necessary because the conditions at the bottom of the loop will never be true if a session is irregular and has to be optimized
		if optimized {
			break
		}

		if tt.totalWork == sessions.duration {
			break
		}

		if sessionEnd.Equal(end) {
			break
		}

		sessionStart = sessionEnd.Add(pauseDurations[i%len(pauseDurations)])
	}

	if len(tt.sessions) == 0 {
		return nil, errors.New("no sessions generated")
	}

	tt.totalDur = tt.sessions[len(tt.sessions)-1].end.Sub(start)

	return &tt, nil
}
