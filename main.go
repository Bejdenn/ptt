package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const (
	valueDelimiter = " "
)

type durationArrayFlag []time.Duration

func (d *durationArrayFlag) String() string {
	return "arrayFlag"
}

func (d *durationArrayFlag) Set(s string) error {
	for _, v := range strings.Split(s, valueDelimiter) {
		item, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("could not parse duration array: %v", err)
		}
		*d = append(*d, item)
	}

	return nil
}

type timeFlag time.Time

func (t *timeFlag) String() string {
	return "time"
}

func (t *timeFlag) Set(s string) error {
	if s == "" {
		return nil
	}

	parsed, err := time.Parse("15:04", s)
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	*t = timeFlag(parsed)
	return nil
}

func main() {
	durationFlag := flag.Duration("duration", time.Duration(0), "Set the working duration that should be covered by pomodoro sessions.")
	sessionLengthFlag := flag.Duration("session-length", 90*time.Minute, "Set the length of a single pomodoro session.")

	var pausePatternFlag durationArrayFlag
	flag.Var(&pausePatternFlag, "pause-pattern", "Set the pause pattern for the pauses between pomodoro sessions. Will be repeated if it has less elements as --duration defines.")

	var startFlag timeFlag
	flag.Var(&startFlag, "start", "Start time of the time table.")

	var endFlag timeFlag
	flag.Var(&endFlag, "end", "Maximum end time of the time table. Ignored if not defined.")

	flag.Parse()

	if (time.Time)(startFlag).IsZero() {
		startFlag = timeFlag(time.Now())
	}
	startFlag = timeFlag(normalize(time.Time(startFlag)))

	tt, err := generateTimetable((time.Time)(startFlag), (time.Time)(endFlag), pausePatternFlag, sessionInfo{*durationFlag, *sessionLengthFlag})
	if err != nil {
		fmt.Printf("could not generate timetable: %v\n", err)
		return
	}

	for _, u := range tt.sessions {
		fmt.Printf("(%d)\t%s\t%s\t%s\n", u.id, u.start.Format(time.TimeOnly), u.end.Format(time.TimeOnly), u.end.Sub(u.start).String())
	}

	fmt.Println("")
	fmt.Printf("Total duration of session: %s\n", tt.totalDur.String())
	fmt.Printf("Total work time: %s\n", tt.totalWork.String())
}

// normalize normalizes a time. Normalization here means to reduce it to its hours and minutes, leaving the rest as
// values of the zero time.
func normalize(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
}
