package timetable

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTimeRange(t *testing.T) {
	type args struct {
		start         time.Time
		end           time.Time
		pause         time.Duration
		duration      time.Duration
		sessionLength time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    TimeRange
		wantErr bool
	}{
		{
			name: "no end, no duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 15, 45, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "no end, duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				duration:      7*time.Hour + 30*time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "end, no duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "end, duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				duration:      3 * time.Hour,
				sessionLength: 90 * time.Minute,
			},
			want: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 12, 15, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimeRange(tt.args.start, tt.args.end, tt.args.pause, tt.args.duration, tt.args.sessionLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTimeRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTimeRangeByDuration(t *testing.T) {
	type args struct {
		start         time.Time
		pause         time.Duration
		duration      time.Duration
		sessionLength time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    TimeRange
		wantErr bool
	}{
		{
			name: "no start, duration",
			args: args{
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
				duration:      7*time.Hour + 30*time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "start, no duration",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "no session length",
			args: args{
				start:    time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:    15 * time.Minute,
				duration: 7*time.Hour + 30*time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "everything ok",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
				duration:      7*time.Hour + 30*time.Minute,
			},
			want: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTimeRangeByDuration(tt.args.start, tt.args.pause, tt.args.duration, tt.args.sessionLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTimeRangeByDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTimeRangeByDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTimeRangeByEnd(t *testing.T) {
	type args struct {
		start         time.Time
		end           time.Time
		pause         time.Duration
		sessionLength time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    TimeRange
		wantErr bool
	}{
		{
			name: "no start",
			args: args{
				end:           time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "no end",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "start after end",
			args: args{
				start:         time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "no session length",
			args: args{
				start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				end:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				pause: 15 * time.Minute,
			},
			want:    TimeRange{},
			wantErr: true,
		},
		{
			name: "everything ok",
			args: args{
				start:         time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				end:           time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTimeRangeByEnd(tt.args.start, tt.args.end, tt.args.pause, tt.args.sessionLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTimeRangeByEnd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTimeRangeByEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTimetable(t *testing.T) {
	type args struct {
		tr            TimeRange
		pause         time.Duration
		sessionLength time.Duration
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
				tr: TimeRange{
					Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
				},
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want: &Timetable{
				Sessions: []Session{
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
							End:   time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no sessions generated",
			args: args{
				tr: TimeRange{
					Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 9, 5, 0, 0, time.UTC),
				},
				pause:         15 * time.Minute,
				sessionLength: 90 * time.Minute,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTimetable(tt.args.tr, tt.args.pause, tt.args.sessionLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTimetable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateTimetable() = %v, want %v", got, tt.want)
			}
		})
	}
}
