package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Bejdenn/ptt/internal/timetable"
)

const (
	TimeOnlyNoSeconds = "15:04"
)

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

		start, err := time.Parse(TimeOnlyNoSeconds, parts[0])
		if err != nil {
			return fmt.Errorf("could not parse start time: %v", err)
		}

		end, err := time.Parse(TimeOnlyNoSeconds, parts[1])
		if err != nil {
			return fmt.Errorf("could not parse end time: %v", err)
		}

		*e = append(*e, timetable.TimeRange{Start: start, End: end})
	}
	return nil
}

type timeFlag time.Time

func (t *timeFlag) String() string {
	return fmt.Sprint(*t)
}

func (t *timeFlag) Set(s string) error {
	parsed, err := time.Parse(TimeOnlyNoSeconds, s)
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	*t = timeFlag(parsed)
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

	if (time.Time)(startFlag).IsZero() {
		startFlag = timeFlag(time.Now())
	}
	startFlag = timeFlag(normalize(time.Time(startFlag)))

	sessions, err := timetable.Generate((time.Time)(startFlag), (time.Time)(endFlag), pauseFlag, durationFlag, sessionLengthFlag, excludesFlag)
	if err != nil {
		fmt.Printf("could not generate timetable: %v\n", err)
		return
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ID\tStart\tEnd\tDuration\tPause\tCumulated Work\tCumulated Time")

	cumulatedWork := time.Duration(0)
	cumulatedTime := time.Duration(0)
	for _, u := range sessions {
		cumulatedWork += u.TimeRange.Duration()
		cumulatedTime += u.TimeRange.Duration() + u.Pause
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t\n", u.ID, u.TimeRange.Start.Format(TimeOnlyNoSeconds), u.TimeRange.End.Format(TimeOnlyNoSeconds), u.TimeRange.Duration(), u.Pause, cumulatedWork, cumulatedTime)
	}

	w.Flush()
}

// normalize normalizes a time. Normalization here means to reduce it to its hours and minutes, leaving the rest as
// values of the zero time.
func normalize(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
}
