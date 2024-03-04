package timetable

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	type args struct {
		start         time.Time
		end           time.Time
		pause         time.Duration
		duration      time.Duration
		sessionLength time.Duration
		excludes      []TimeRange
	}
	tests := []struct {
		name    string
		args    args
		want    SessionSlice
		wantErr bool
	}{
		{
			name: "end, duration",
			args: args{
				start:         time.Date(2020, 1, 1, 7, 30, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 19, 30, 0, 0, time.UTC),
				pause:         10 * time.Minute,
				duration:      6 * time.Hour,
				sessionLength: 90 * time.Minute,
				excludes: []TimeRange{
					{
						Start: time.Date(2020, 1, 1, 8, 30, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 10, 15, 0, 0, time.UTC),
					},
					{
						Start: time.Date(2020, 1, 1, 12, 45, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 16, 15, 0, 0, time.UTC),
					},
				},
			},
			want: []Session{
				{
					ID: 1,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 7, 30, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 8, 30, 0, 0, time.UTC),
					},
				},
				{
					ID: 2,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 10, 15, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 45, 0, 0, time.UTC),
					},
					Pause: 10 * time.Minute,
				},
				{
					ID: 3,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 11, 55, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 12, 45, 0, 0, time.UTC),
					},
				},
				{
					ID: 4,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 16, 15, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 17, 45, 0, 0, time.UTC),
					},
					Pause: 10 * time.Minute,
				},
				{
					ID: 5,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 17, 55, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 19, 5, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "end, no duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 14, 30, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want: []Session{
				{
					ID: 1,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
					},
					Pause: 15 * time.Minute,
				},
				{
					ID: 2,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 10, 45, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 12, 15, 0, 0, time.UTC),
					},
					Pause: 15 * time.Minute,
				},
				{
					ID: 3,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 14, 0, 0, 0, time.UTC),
					},
					Pause: 15 * time.Minute,
				},
				{
					ID: 4,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 14, 15, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 14, 30, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "end, no duration, with excludes",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 14, 30, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
				excludes: []TimeRange{
					{
						Start: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC),
					},
					{
						Start: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []Session{
				{
					ID: 1,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
					},
				},
				{
					ID: 2,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
				{
					ID: 3,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 14, 30, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no end, duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				duration:      4 * time.Hour,
				sessionLength: 90 * time.Minute,
			},
			want: []Session{
				{
					ID: 1,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
					},
					Pause: 15 * time.Minute,
				},
				{
					ID: 2,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 10, 45, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 12, 15, 0, 0, time.UTC),
					},
					Pause: 15 * time.Minute,
				},
				{
					ID: 3,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 13, 30, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no end, duration, with excludes",
			args: args{
				start:         time.Date(2020, 1, 1, 14, 2, 0, 0, time.UTC),
				pause:         10 * time.Minute,
				duration:      6 * time.Hour,
				sessionLength: 90 * time.Minute,
				excludes: []TimeRange{
					{
						Start: time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 16, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []Session{
				{
					ID: 1,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 14, 2, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
					},
				},
				{
					ID: 2,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 16, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
					},
					Pause: 10 * time.Minute,
				},
				{
					ID: 3,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 17, 40, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 19, 10, 0, 0, time.UTC),
					},
					Pause: 10 * time.Minute,
				},
				{
					ID: 4,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 19, 20, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 20, 50, 0, 0, time.UTC),
					},
					Pause: 10 * time.Minute,
				},
				{
					ID: 5,
					TimeRange: TimeRange{
						Start: time.Date(2020, 1, 1, 21, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 21, 32, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.start, tt.args.end, tt.args.pause, tt.args.duration, tt.args.sessionLength, tt.args.excludes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
