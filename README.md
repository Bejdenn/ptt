# ptt - pomodoro time table

## Description

A simple pomodoro timetable for the terminal. It gives you an overview of your pomodoro sessions and their times.

## Installation

Currently, the installation is only possible from source. To install it, you need to have [go](https://golang.org/)
installed.

```bash
go install
```

This will install the binary in your `$GOPATH/bin` directory. Make sure that this directory is in your `$PATH`.

## Usage

```
Usage:
    ptt [--start] [--end] [--session-length] [--duration] [--pause-pattern]

Options:
    --start             Start time of the time table. Default is current time.
    --end               Maximum end time of the time table. Ignored if not defined.
    --session-length    Set the length of a single pomodoro session. Default is 90 minutes.
    --duration          Set the working duration that should be covered by pomodoro sessions.
    --pause-pattern     Set the pause pattern for the pauses between pomodoro sessions. Will be repeated if it has less elements as --duration defines.
```

For example, you want to work 6 hours with 90 minutes sessions and 15 minutes breaks. You can use the following command:

```bash
ptt --duration 6h --session-length 90m --pause-pattern 15m
```

```
# Output (assuming the current time is 8:00):
(1)	08:00:00	09:30:00	1h30m0s
(2)	09:45:00	11:15:00	1h30m0s
(3)	11:30:00	13:00:00	1h30m0s
(4)	13:15:00	14:45:00	1h30m0s

Total duration of session: 6h45m0s
Total work time: 6h0m0s
```
