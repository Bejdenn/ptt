package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	durationFlag := flag.Duration("duration", time.Duration(0), "Set the working duration that should be covered by pomodoro sessions.")
	sessionLengthFlag := flag.Duration("session-length", 90*time.Minute, "Set the length of a single pomodoro session.")
	pausePatternFlag := flag.String("pause-pattern", "10m", "Set the pause pattern for the pauses between pomodoro sessions. Will be repeated if it has less elements as --duration defines.")

	startFlag := flag.String("start", time.Now().Format("15:04"), "Start time of the time table.")
	endFlag := flag.String("end", "", "Maximum end time of the time table. Ignored if not defined.")

	flag.Parse()

	start, err := time.Parse("15:04", *startFlag)
	if err != nil {
		fmt.Printf("could not parse start time %v\n", err)
		return
	}

	var end time.Time
	if *endFlag != "" {
		end, err = time.Parse("15:04", *endFlag)
		if err != nil {
			fmt.Printf("could not parse end time: %v\n", err)
			return
		}
	}

	if pausePatternFlag == nil {
		*pausePatternFlag = ""
	}

	tt, err := generateTimetable(start, end, *pausePatternFlag, sessionInfo{*durationFlag, *sessionLengthFlag})
	if err != nil {
		fmt.Printf("could not generate timetable: %v\n", err)
		return
	}

	for _, u := range tt.sessions {
		fmt.Printf("(%d)\t%s\t%s\t%s\n", u.id, u.start.Format("15:04"), u.end.Format("15:04"), u.end.Sub(u.start).String())
	}

	fmt.Println("")
	fmt.Printf("Total duration of session: %s\n", tt.totalDur.String())
	fmt.Printf("Total work time: %s\n", tt.totalWork.String())
}
