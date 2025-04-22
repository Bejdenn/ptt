package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"

	"time"

	"github.com/Bejdenn/ptt/internal/timetable"
	"github.com/Bejdenn/timerange"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config[T any] struct {
	Defaults struct {
		Duration         T `toml:"duration" env-default:"-1ns"`
		MinSessionLength T `toml:"min-session-length" env-default:"15m"`
		MaxSessionLength T `toml:"max-session-length" env-default:"90m"`
		Pause            T `toml:"pause" env-default:"15m"`
	} `toml:"defaults"`
}

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

const usage = `Usage:
    ptt [-s START] [-e END] [-l LENGTH] [-L LENGTH] [-d DURATION] [-p PAUSE] (-x EXCLUDE)...
Options:
    -s, --start START                Set START as the start time of the time table. Default is current time.
    -e, --end END                    Set END as the end time of the time table. Ignored if not defined.
    -l, --min-session-length LENGTH  Set LENGTH as the minimum length of a single pomodoro session. Default is 15 minutes.
    -L, --max-session-length LENGTH  Set LENGTH as the maximum length of a single pomodoro session. Default is 90 minutes.
    -d, --duration DURATION          Set DURATION as the working duration that should be covered by pomodoro sessions.
    -p, --pause PAUSE                Set PAUSE as the pause duration between pomodoro sessions.
    -x, --exclude EXCLUDE            Exclude EXCLUDE to prevent from being overlapped by a pomodoro session. Can be repeated.
	
END and DURATION are mutually exclusive. If both are defined, the time table will used that ends earlier.
Defining no END or DURATION will result in an error.
If any of the END or DURATION yield an empty slice of sessions, no sessions will be generated.

The defaults can be set globally with a configuration file at ($HOME/.config/ptt/config.toml).

The format of the durations and time values are the same that the Go programming language uses for its time parsing.`

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
}

func durationParseErr(err error) error {
	return fmt.Errorf("error while parsing duration: %v\n", err)
}

func ReadConfig() Config[time.Duration] {
	var cfg Config[string]

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while getting user config directory: %v\n", err)
		os.Exit(1)
	}

	err = cleanenv.ReadConfig(filepath.Join(home, ".config", "ptt", "config.toml"), &cfg)
	if err != nil {
		fmt.Printf("error while reading config: %v\n", err)
		os.Exit(1)
	}

	duration, err := time.ParseDuration(cfg.Defaults.Duration)
	if err != nil {
		fmt.Fprint(os.Stderr, durationParseErr(err))
		os.Exit(1)
	}

	minSessionLength, err := time.ParseDuration(cfg.Defaults.MinSessionLength)
	if err != nil {
		fmt.Fprint(os.Stderr, durationParseErr(err))
		os.Exit(1)
	}

	maxSessionLength, err := time.ParseDuration(cfg.Defaults.MaxSessionLength)
	if err != nil {
		fmt.Fprint(os.Stderr, durationParseErr(err))
		os.Exit(1)
	}

	pause, err := time.ParseDuration(cfg.Defaults.Pause)
	if err != nil {
		fmt.Fprint(os.Stderr, durationParseErr(err))
		os.Exit(1)
	}

	return Config[time.Duration]{Defaults: struct {
		Duration         time.Duration "toml:\"duration\" env-default:\"-1ns\""
		MinSessionLength time.Duration "toml:\"min-session-length\" env-default:\"15m\""
		MaxSessionLength time.Duration "toml:\"max-session-length\" env-default:\"90m\""
		Pause            time.Duration "toml:\"pause\" env-default:\"15m\""
	}{
		duration, minSessionLength, maxSessionLength, pause,
	}}
}

func main() {
	cfg := ReadConfig()

	var (
		startFlag            = NewTimeFlag(now)
		endFlag              = NewTimeFlag(time.Time{})
		minSessionLengthFlag time.Duration
		maxSessionLengthFlag time.Duration
		durationFlag         time.Duration
		pauseFlag            time.Duration
		excludesFlag         excludesMultiFlag
		versionFlag          bool
	)

	flag.Var(&startFlag, "start", "set the start time")
	flag.Var(&startFlag, "s", "set the start time (shorthand)")

	flag.Var(&endFlag, "end", "set the end time")
	flag.Var(&endFlag, "e", "set the end time (shorthand)")

	flag.DurationVar(&minSessionLengthFlag, "min-session-length", cfg.Defaults.MinSessionLength, "set the session length")
	flag.DurationVar(&minSessionLengthFlag, "l", cfg.Defaults.MinSessionLength, "set the session length (shorthand)")

	flag.DurationVar(&maxSessionLengthFlag, "max-session-length", cfg.Defaults.MaxSessionLength, "set the session length")
	flag.DurationVar(&maxSessionLengthFlag, "L", cfg.Defaults.MaxSessionLength, "set the session length (shorthand)")

	flag.DurationVar(&durationFlag, "duration", cfg.Defaults.Duration, "set the working duration")
	flag.DurationVar(&durationFlag, "d", cfg.Defaults.Duration, "set the working duration (shorthand)")

	flag.DurationVar(&pauseFlag, "pause", cfg.Defaults.Pause, "set the pause duration")
	flag.DurationVar(&pauseFlag, "p", cfg.Defaults.Pause, "set the pause duration (shorthand)")

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

	sessionLength, err := timetable.NewSessionLength(minSessionLengthFlag, maxSessionLengthFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot generate timetable: %v\n", err)
		os.Exit(1)
	}

	sessions, err := timetable.Generate(start, end, pauseFlag, durationFlag, sessionLength, excludesFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot generate timetable: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(sessions)
}
