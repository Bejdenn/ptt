package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"

	"time"

	"github.com/Bejdenn/ptt/internal/timetable"
	"github.com/Bejdenn/timerange"
)

var now = time.Now().Truncate(time.Minute)

type excludesMultiFlag []timerange.TimeRange

func (e *excludesMultiFlag) String() string {
	return fmt.Sprint(*e)
}

var rangePattern = regexp.MustCompile(`\d{2}:\d{2}-\d{2}:\d{2}`)

func (e *excludesMultiFlag) Set(value string) error {
	exclude, err := timerange.Parse(now, rangePattern.FindAllString(value, -1))
	if err != nil {
		return err
	}
	*e = append(*e, exclude...)
	return nil
}

type timeFlag time.Time

func NewTimeFlag(t time.Time) timeFlag {
	// avoid altering zero time, as we indicate absence of end time that way (see https://github.com/Bejdenn/ptt/issues/28)
	if !t.IsZero() {
		t = timerange.Normalize(now, t)
	}
	return timeFlag(t)
}

func (t *timeFlag) String() string {
	return fmt.Sprint(*t)
}

func (t *timeFlag) Set(s string) error {
	parsed, err := time.Parse(timerange.TimeOnlyNoSeconds, s)
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	*t = NewTimeFlag(parsed)
	return nil
}

const (
	defaultDuration      = time.Duration(-1)
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
Defining no END or DURATION will result in an error.
If any of the END or DURATION yield an empty slice of sessions, no sessions will be generated.

The format of the durations and time values are the same that the Go programming language uses for its time parsing.`

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
}

func main() {
	var (
		startFlag         = NewTimeFlag(now)
		endFlag           = NewTimeFlag(time.Time{})
		sessionLengthFlag time.Duration
		durationFlag      time.Duration
		pauseFlag         time.Duration
		excludesFlag      excludesMultiFlag
		versionFlag       bool
	)

	flag.Var(&startFlag, "start", "set the start time")
	flag.Var(&startFlag, "s", "set the start time (shorthand)")

	flag.Var(&endFlag, "end", "set the end time")
	flag.Var(&endFlag, "e", "set the end time (shorthand)")

	flag.DurationVar(&sessionLengthFlag, "session-length", defaultSessionLength, "set the session length")
	flag.DurationVar(&sessionLengthFlag, "l", defaultSessionLength, "set the session length (shorthand)")

	flag.DurationVar(&durationFlag, "duration", defaultDuration, "set the working duration")
	flag.DurationVar(&durationFlag, "d", defaultDuration, "set the working duration (shorthand)")

	flag.DurationVar(&pauseFlag, "pause", defaultPause, "set the pause duration")
	flag.DurationVar(&pauseFlag, "p", defaultPause, "set the pause duration (shorthand)")

	flag.Var(&excludesFlag, "exclude", "exclude one or several time ranges")
	flag.Var(&excludesFlag, "x", "exclude one or several time ranges (shorthand)")
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

	start, end := time.Time(startFlag), time.Time(endFlag)

	if !end.IsZero() && end.Before(start) {
		// avoid altering zero time, as we indicate absence of end time that way (see https://github.com/Bejdenn/ptt/issues/28)
		end = end.AddDate(0, 0, 1)
	}

	sessions, err := timetable.Generate(start, end, pauseFlag, durationFlag, sessionLengthFlag, excludesFlag)
	if err != nil {
		fmt.Printf("cannot generate timetable: %v\n", err)
		return
	}

	fmt.Print(sessions)
}
