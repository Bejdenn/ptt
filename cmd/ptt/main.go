package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Bejdenn/ptt/internal/timetable"
)

var now = time.Now()

type excludesMultiFlag []timetable.TimeRange

func (e *excludesMultiFlag) String() string {
	return fmt.Sprint(*e)
}

func (e *excludesMultiFlag) Set(s string) error {
	excludes := strings.Split(s, " ")
	for _, exclude := range excludes {
		parts := strings.Split(exclude, "-")
		if len(parts) != 2 {
			return fmt.Errorf("could not parse exclude: %s", exclude)
		}

		start, err := time.Parse(timetable.TimeOnlyNoSeconds, parts[0])
		if err != nil {
			return fmt.Errorf("could not parse start time: %v", err)
		}

		end, err := time.Parse(timetable.TimeOnlyNoSeconds, parts[1])
		if err != nil {
			return fmt.Errorf("could not parse end time: %v", err)
		}

		*e = append(*e, timetable.TimeRange{Start: normalize(start), End: normalize(end)})
	}
	return nil
}

type timeFlag time.Time

func (t *timeFlag) String() string {
	return fmt.Sprint(*t)
}

func (t *timeFlag) Set(s string) error {
	parsed, err := time.Parse(timetable.TimeOnlyNoSeconds, s)
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	*t = timeFlag(normalize(parsed))
	return nil
}

const (
	defaultDuration      = time.Duration(0)
	defaultSessionLength = 90 * time.Minute
	defaultPause         = 15 * time.Minute
)

const usage = `Usage:
    ptt [-s START] [-e END] [-l LENGTH] [-d DURATION] [-p PAUSE] (-x EXCLUDE)...

Options:
    -s, --start START            Set START as the start time of the time table. Default is current time.
    -e, --end END                Set END as the end time of the time table. Ignored if not defined.
    -l, --session-length LENGTH  Set LENGTH as the length of a single pomodoro session. Default is 90 minutes.
    -d, --duration DURATION      Set DURATION as the working duration that should be covered by pomodoro sessions.
    -p, --pause PAUSE            Set PAUSE as the pause duration between pomodoro sessions.
    -x, --exclude EXCLUDE        Exclude EXCLUDE to prevent from being overlapped by a pomodoro session. Can be repeated.
	
END and DURATION are mutually exclusive. If both are defined, the time table will used that ends earlier.
The format of the durations and time values can be be set as the Go programming language's parsing format defines it.`

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
}

func main() {
	var (
		durationFlag      time.Duration
		sessionLengthFlag time.Duration
		pauseFlag         time.Duration
		startFlag         timeFlag
		endFlag           timeFlag
		excludesFlag      excludesMultiFlag
		versionFlag       bool
	)

	flag.DurationVar(&durationFlag, "duration", time.Duration(0), "set the working duration")
	flag.DurationVar(&durationFlag, "d", defaultDuration, "set the working duration (shorthand)")
	flag.DurationVar(&sessionLengthFlag, "session-length", defaultSessionLength, "set the session length")
	flag.DurationVar(&sessionLengthFlag, "sl", defaultSessionLength, "set the session length (shorthand)")
	flag.DurationVar(&pauseFlag, "pause", defaultPause, "set the pause duration")
	flag.DurationVar(&pauseFlag, "p", defaultPause, "set the pause duration (shorthand)")
	flag.Var(&startFlag, "start", "set the start time")
	flag.Var(&startFlag, "s", "set the start time (shorthand)")
	flag.Var(&endFlag, "end", "set the end time")
	flag.Var(&endFlag, "e", "set the end time (shorthand)")
	flag.Var(&excludesFlag, "exclude", "exclude one or several time ranges")
	flag.Var(&excludesFlag, "ex", "exclude one or several time ranges (shorthand)")
	flag.BoolVar(&versionFlag, "version", false, "Print the version and exit.")
	flag.Parse()

	if versionFlag {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(buildInfo.Main.Version)
			return
		}
		fmt.Println("(unknown)")
		return
	}

	start, end := (time.Time)(startFlag), (time.Time)(endFlag)

	if start.IsZero() {
		start = now
	}

	// if end is before start, interpret it as being the same time on the next day
	if !end.IsZero() && end.Before(start) {
		end = end.AddDate(0, 0, 1)
	}

	sessions, err := timetable.Generate(start, end, pauseFlag, durationFlag, sessionLengthFlag, excludesFlag)
	if err != nil {
		fmt.Printf("cannot generate timetable: %v\n", err)
		return
	}

	fmt.Print(sessions)
}

func normalize(t time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
}
