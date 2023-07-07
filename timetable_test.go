package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_generateTimetable(t *testing.T) {
	type args struct {
		start        time.Time
		end          time.Time
		pausePattern string
		sessions     sessionInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *timetable
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				start:        time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				pausePattern: "10",
				sessions:     sessionInfo{duration: time.Duration(2) * time.Hour, sessionLength: 30},
			},
			want: &timetable{
				sessions: []session{
					{
						id:    1,
						start: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
					},
					{
						id:    2,
						start: time.Date(2020, 1, 1, 10, 40, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 11, 10, 0, 0, time.UTC),
					},
					{
						id:    3,
						start: time.Date(2020, 1, 1, 11, 20, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 11, 50, 0, 0, time.UTC),
					},
					{
						id:    4,
						start: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
					},
				},
				totalWork: time.Duration(120) * time.Minute,
				totalDur:  time.Duration(150) * time.Minute,
			},
		}, {
			name: "different pause times",
			args: args{
				start:        time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				pausePattern: "10-5",
				sessions:     sessionInfo{duration: time.Duration(2) * time.Hour, sessionLength: 30},
			},
			want: &timetable{
				sessions: []session{
					{
						id:    1,
						start: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
					},
					{
						id:    2,
						start: time.Date(2020, 1, 1, 10, 40, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 11, 10, 0, 0, time.UTC),
					},
					{
						id:    3,
						start: time.Date(2020, 1, 1, 11, 15, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 11, 45, 0, 0, time.UTC),
					},
					{
						id:    4,
						start: time.Date(2020, 1, 1, 11, 55, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 12, 25, 0, 0, time.UTC),
					},
				},
				totalWork: time.Duration(120) * time.Minute,
				totalDur:  time.Duration(145) * time.Minute,
			},
		}, {
			name: "only end time",
			args: args{
				start:        time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
				end:          time.Date(2020, 1, 1, 20, 45, 0, 0, time.UTC),
				pausePattern: "10",
				sessions:     sessionInfo{duration: time.Duration(0), sessionLength: 90},
			},
			want: &timetable{
				sessions: []session{
					{
						id:    1,
						start: time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 20, 15, 0, 0, time.UTC),
					},
					{
						id:    2,
						start: time.Date(2020, 1, 1, 20, 25, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 20, 45, 0, 0, time.UTC),
					},
				},
				totalWork: time.Duration(110) * time.Minute,
				totalDur:  time.Duration(120) * time.Minute,
			},
		},
		{
			name: "session length not multiple of duration",
			args: args{
				start:        time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
				pausePattern: "10",
				sessions:     sessionInfo{duration: time.Duration(2) * time.Hour, sessionLength: 45},
			},
			want: &timetable{
				sessions: []session{
					{
						id:    1,
						start: time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 19, 30, 0, 0, time.UTC),
					},
					{
						id:    2,
						start: time.Date(2020, 1, 1, 19, 40, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 20, 25, 0, 0, time.UTC),
					},
					{
						id:    3,
						start: time.Date(2020, 1, 1, 20, 35, 0, 0, time.UTC),
						end:   time.Date(2020, 1, 1, 21, 5, 0, 0, time.UTC),
					},
				},
				totalWork: time.Duration(120) * time.Minute,
				totalDur:  time.Duration(140) * time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateTimetable(tt.args.start, tt.args.end, tt.args.pausePattern, tt.args.sessions)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateTimetable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateTimetable() = %v, want %v", got, tt.want)
			}
		})
	}
}
