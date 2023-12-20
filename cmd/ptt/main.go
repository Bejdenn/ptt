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

	flag.DurationVar(&durationFlag, "duration", time.Duration(0), "Set the working duration that should be covered by pomodoro sessions.")
	flag.DurationVar(&sessionLengthFlag, "session-length", 90*time.Minute, "Set the length of a single pomodoro session.")
	flag.DurationVar(&pauseFlag, "pause", 15*time.Minute, "Set the duration for the pauses between pomodoro sessions.")
	flag.Var(&startFlag, "start", "Start time of the time table.")
	flag.Var(&endFlag, "end", "Maximum end time of the time table. Ignored if not defined.")
	flag.Var(&excludesFlag, "exclude", "Exclude multiple time ranges that should not be covered by pomodoro sessions.")
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
