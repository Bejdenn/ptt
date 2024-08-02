package timerange

import (
	"reflect"
	"testing"
	"time"
)

func Test_max(t *testing.T) {
	type args struct {
		t time.Time
		u time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t after u",
			args: args{
				t: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				u: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			name: "t before u",
			args: args{
				t: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
				u: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			name: "t equal u",
			args: args{
				t: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				u: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.t, tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		t time.Time
		u time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t after u",
			args: args{
				t: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				u: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
		},
		{
			name: "t before u",
			args: args{
				t: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
				u: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
		},
		{
			name: "t equal u",
			args: args{
				t: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				u: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.t, tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TimeRange_sub(t *testing.T) {
	type args struct {
		u TimeRange
	}
	tests := []struct {
		name    string
		tr      TimeRange
		args    args
		wantRes []TimeRange
	}{
		{
			name: "u.Start before tr.Start, u.End before tr.Start",
			tr: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			args: args{
				u: TimeRange{
					Start: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				},
			},
			wantRes: []TimeRange{
				{
					Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "u.Start before tr.Start, u.End after tr.End",
			tr: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			args: args{
				u: TimeRange{
					Start: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC),
				},
			},
			wantRes: []TimeRange{},
		},
		{
			name: "u in between tr",
			tr: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			args: args{
				u: TimeRange{
					Start: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 17, 0, 0, 0, time.UTC),
				},
			},
			wantRes: []TimeRange{
				{
					Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2020, 1, 1, 17, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.tr.sub(tt.args.u); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("TimeRange.sub() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_TimeRange_subMulti(t *testing.T) {
	type args struct {
		us []TimeRange
	}
	tests := []struct {
		name string
		tr   TimeRange
		args args
		want []TimeRange
	}{
		{
			name: "three excludes",
			tr: TimeRange{
				Start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 1, 1, 17, 30, 0, 0, time.UTC),
			},
			args: args{
				us: []TimeRange{
					{
						Start: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC),
					},
					{
						Start: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
					},
					{
						Start: time.Date(2020, 1, 1, 14, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []TimeRange{
				{
					Start: time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
					End:   time.Date(2020, 1, 1, 14, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.SubMulti(tt.args.us); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeRange.subMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTimeRange(t *testing.T) {
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
			got, err := New(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTimeRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
