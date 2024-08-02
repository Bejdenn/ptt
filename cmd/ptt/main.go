package main

import (
	"fmt"
	"os"

	"strings"
	"time"

	"github.com/Bejdenn/ptt/internal/timetable"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
END and DURATION are mutually exclusive. If both are defined, the time table will use the value that results in the earlier end time.
The format of the durations and time values are the same that the Go programming language uses for its time parsing.
`, cli.AppHelpTemplate)

	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h"},
		Usage:              "Show help",
		DisableDefaultText: true,
	}

	(&cli.App{
		Name:  "ptt",
		Usage: "Pomodoro time table for the terminal. Interpolates your working times for a better daily overview.",
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name:        "start",
				Layout:      timetable.TimeOnlyNoSeconds,
				Aliases:     []string{"s"},
				Usage:       "Set `START` as the start time of the time table.",
				DefaultText: "current time",
				Value:       cli.NewTimestamp(time.Now()),
			},
			&cli.TimestampFlag{
				Name:    "end",
				Layout:  timetable.TimeOnlyNoSeconds,
				Aliases: []string{"e"},
				Usage:   "Set `END` as the end time of the time table. Ignored if not defined.",
				Value:   cli.NewTimestamp(time.Time{}),
				Action: func(c *cli.Context, end *time.Time) error {
					if !end.IsZero() && end.Before(*c.Timestamp("start")) {
						*end = end.AddDate(0, 0, 1)
					}
					return nil
				},
			},
			&cli.DurationFlag{
				Name:    "session-length",
				Value:   90 * time.Minute,
				Aliases: []string{"l"},
				Usage:   "Set `LENGTH` as the length of a single pomodoro session.",
			},
			&cli.DurationFlag{
				Name:    "duration",
				Value:   time.Duration(0),
				Aliases: []string{"d"},
				Usage:   "Set `DURATION` as the working duration that should be covered by pomodoro sessions.",
			},
			&cli.DurationFlag{
				Name:    "pause",
				Value:   15 * time.Minute,
				Aliases: []string{"p"},
				Usage:   "Set `PAUSE` as the pause duration between pomodoro sessions.",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"x"},
				Usage:   "Exclude `EXCLUDE` to prevent from being overlapped by a pomodoro session. Can be repeated.",
			},
		},
		Action: func(c *cli.Context) error {
			excludes, err := parseTimeRanges(*c.Timestamp("start"), c.StringSlice("exclude"))
			if err != nil {
				return fmt.Errorf("cannot parse excludes: %v\n", err)
			}

			sessions, err := timetable.Generate(*c.Timestamp("start"), *c.Timestamp("end"), c.Duration("pause"), c.Duration("duration"), c.Duration("session-length"), excludes)
			if err != nil {
				return fmt.Errorf("cannot generate timetable: %v\n", err)
			}

			fmt.Print(sessions)
			return nil
		},
		HideHelpCommand: true,
	}).Run(os.Args)
}

func parseTimeRanges(ref time.Time, args []string) ([]timetable.TimeRange, error) {
	excludes := make([]timetable.TimeRange, 0, len(args))
	for _, exclude := range args {
		parts := strings.Split(exclude, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("exclude did not contain hyphen delimiter: %s", exclude)
		}

		start, err := time.Parse(timetable.TimeOnlyNoSeconds, parts[0])
		if err != nil {
			return nil, fmt.Errorf("could not parse start time: %v", err)
		}

		end, err := time.Parse(timetable.TimeOnlyNoSeconds, parts[1])
		if err != nil {
			return nil, fmt.Errorf("could not parse end time: %v", err)
		}

		excludes = append(excludes, timetable.TimeRange{Start: normalize(ref, start), End: normalize(ref, end)})
	}
	return excludes, nil
}

func normalize(ref, t time.Time) time.Time {
	return time.Date(ref.Year(), ref.Month(), ref.Day(), t.Hour(), t.Minute(), 0, 0, ref.Location())
}
