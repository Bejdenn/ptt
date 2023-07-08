# ptt - pomodoro time table

## Description

A simple pomodoro time table for the terminal. It gives you an overview of your pomodoro sessions and their times.

## Installation

Currently, the installation is only possible from source. To install it, you need to have [go](https://golang.org/) installed.

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

For example, you want to work 6 hours with 25 minutes sessions and 5 minutes breaks. You can use the following command:

```bash
ptt --duration 4h --session-length 25 --pause-pattern 5

# Output (assuming the current time is 8:00):
(1)     08:00   08:25   25m0s
(2)     08:30   08:55   25m0s
(3)     09:00   09:25   25m0s
(4)     09:30   09:55   25m0s
(5)     10:00   10:25   25m0s
(6)     10:30   10:55   25m0s
(7)     11:00   11:25   25m0s
(8)     11:30   11:55   25m0s
(9)     12:00   12:25   25m0s
(10)    12:30   12:45   15m0s

Total duration of session: 4h45m0s
Total work time: 4h0m0s
```
