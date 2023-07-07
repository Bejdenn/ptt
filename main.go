package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	step = 5
)

func main() {
	durationFlag := flag.String("duration", "", "number of units to learn")
	sessionLengthFlag := flag.Int("session-length", 90, "length of each unit in minutes")
	pausePatternFlag := flag.String("pause-pattern", "10", "pattern that specifies how long the pause after each unit should be")

	startFlag := flag.String("start", time.Now().Format("15:04"), "string of start time")
	endFlag := flag.String("end", "", "string of end time")

	flag.Parse()

	var duration time.Duration
	if *durationFlag != "" {
		var err error
		duration, err = time.ParseDuration(*durationFlag)
		if err != nil {
			fmt.Printf("could not parse duration: %v\n", err)
			return
		}
	}

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

	tt, err := generateTimetable(start, end, *pausePatternFlag, sessionInfo{duration, *sessionLengthFlag})
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
