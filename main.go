package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	valueDelimiter    = " "
	TimeOnlyNoSeconds = "15:04"
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

	parsed, err := time.Parse(TimeOnlyNoSeconds, s)
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

	var versionFlag bool
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

	tt, err := generateTimetable((time.Time)(startFlag), (time.Time)(endFlag), pausePatternFlag, sessionInfo{*durationFlag, *sessionLengthFlag})
	if err != nil {
		fmt.Printf("could not generate timetable: %v\n", err)
		return
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ID\tStart\tEnd\tDuration\tPause\tCumulated Work\tCumulated Time")

	cumulatedWork := time.Duration(0)
	cumulatedTime := time.Duration(0)
	for _, u := range tt.sessions {
		cumulatedWork += u.end.Sub(u.start)
		cumulatedTime += u.end.Sub(u.start) + u.pause
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t\n", u.id, u.start.Format(TimeOnlyNoSeconds), u.end.Format(TimeOnlyNoSeconds), u.end.Sub(u.start), u.pause, cumulatedWork, cumulatedTime)
	}

	w.Flush()
}

// normalize normalizes a time. Normalization here means to reduce it to its hours and minutes, leaving the rest as
// values of the zero time.
func normalize(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
}
