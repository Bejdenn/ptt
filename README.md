# ptt - pomodoro time table

## Description

A simple pomodoro timetable for the terminal. It gives you an overview of your pomodoro sessions and their times.

## Installation

Currently, the installation is only possible from source. To install it, you need to have [go](https://golang.org/)
installed.

```bash
go install github.com/Bejdenn/ptt/cmd/...@latest
```

This will install the binary in your `$GOPATH/bin` directory. Make sure that this directory is in your `$PATH`.

## Usage

```
Usage:
    ptt [-s START] [-e END] [-l LENGTH] [-d DURATION] [-p PAUSE] (-x EXCLUDE)...

Options:
    -s, --start START            Set START as the start time of the time table. Default is current time.
    -e, --end END                Set END as the end time of the time table. Ignored if not defined.
    -l, --session-length LENGTH  Set LENGTH as the length of a single pomodoro session. Default is 90 minutes.
    -d, --duration DURATION      Set DURATION as the working duration that should be covered by pomodoro sessions.
    -p, --pause PAUSE            Set PAUSE as the pause duration between pomodoro sessions.
    -x, --exclude EXCLUDE        Exclude EXCLUDE to prevent from being overlapped by a pomodoro session. Can be repeated.

END and DURATION are mutually exclusive. If both are defined, the time table will used that ends earlier.
The format of the durations and time values are the same that the Go programming language uses for its time parsing.
```

For example, you want to work 6 hours with 90 minutes sessions and 15 minutes breaks. You can use the following command:

```bash
ptt -d 6h -l 90m -p 15m
```

```
# Output (assuming the current time is 8:00):
ID   Start      End        Duration   Pause   Cumulated Work   Cumulated Time
1    08:00:00   09:30:00   1h30m0s    15m0s   1h30m0s          1h45m0s
2    09:45:00   11:15:00   1h30m0s    15m0s   3h0m0s           3h30m0s
3    11:30:00   13:00:00   1h30m0s    15m0s   4h30m0s          5h15m0s
4    13:15:00   14:45:00   1h30m0s    0s      6h0m0s           6h45m0s
```
