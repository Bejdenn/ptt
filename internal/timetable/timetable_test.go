package timetable

import (
	"reflect"
	"testing"
	"time"
)

func Test_generateTimetable(t *testing.T) {
	type args struct {
		start    time.Time
		end      time.Time
		pause    time.Duration
		sessions SessionInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *Timetable
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				start:    time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				pause:    10 * time.Minute,
				sessions: SessionInfo{Duration: time.Duration(2) * time.Hour, SessionLength: time.Duration(30) * time.Minute},
			},
			want: &Timetable{
				Sessions: []Session{
					{
						ID:    1,
						Start: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    2,
						Start: time.Date(2020, 1, 1, 10, 40, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 10, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    3,
						Start: time.Date(2020, 1, 1, 11, 20, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 50, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    4,
						Start: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
					},
				},
			},
		}, {
			name: "different pause times",
			args: args{
				start:    time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				pause:    10 * time.Minute,
				sessions: SessionInfo{Duration: time.Duration(2) * time.Hour, SessionLength: time.Duration(30) * time.Minute},
			},
			want: &Timetable{
				Sessions: []Session{
					{
						ID:    1,
						Start: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    2,
						Start: time.Date(2020, 1, 1, 10, 40, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 10, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    3,
						Start: time.Date(2020, 1, 1, 11, 20, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 50, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    4,
						Start: time.Date(2020, 1, 1, 12, 00, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
					},
				},
			},
		}, {
			name: "only end time",
			args: args{
				start:    time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
				end:      time.Date(2020, 1, 1, 20, 45, 0, 0, time.UTC),
				pause:    10 * time.Minute,
				sessions: SessionInfo{Duration: time.Duration(0), SessionLength: time.Duration(90) * time.Minute},
			},
			want: &Timetable{
				Sessions: []Session{
					{
						ID:    1,
						Start: time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 20, 15, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    2,
						Start: time.Date(2020, 1, 1, 20, 25, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 20, 45, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "session length not multiple of duration",
			args: args{
				start:    time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
				pause:    10 * time.Minute,
				sessions: SessionInfo{Duration: time.Duration(2) * time.Hour, SessionLength: time.Duration(45) * time.Minute},
			},
			want: &Timetable{
				Sessions: []Session{
					{
						ID:    1,
						Start: time.Date(2020, 1, 1, 18, 45, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 19, 30, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    2,
						Start: time.Date(2020, 1, 1, 19, 40, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 20, 25, 0, 0, time.UTC),
						Pause: time.Duration(10) * time.Minute,
					},
					{
						ID:    3,
						Start: time.Date(2020, 1, 1, 20, 35, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 21, 5, 0, 0, time.UTC),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTimetable(tt.args.start, tt.args.end, tt.args.pause, tt.args.sessions)
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
