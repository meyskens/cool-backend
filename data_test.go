package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseHex(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 255",
			args: args{
				in: "FF",
			},
			wantErr: false,
			want:    "255",
		},
		{
			name: "dummy input",
			args: args{
				in: "21b937",
			},
			wantErr: false,
			want:    "2210103",
		},
		{
			name: "dummy input",
			args: args{
				in: "013461B7",
			},
			wantErr: false,
			want:    "20210103",
		},
		{
			name: "dummy input",
			args: args{
				in: "013ee8a0",
			},
			wantErr: false,
			want:    "20900000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHex(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeDataSane(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test without 00",
			args: args{
				in: "FF10",
			},
			want: "10FF",
		},
		{
			name: "test with 00",
			args: args{
				in: "FF00",
			},
			want: "FF",
		},
		{
			name: "test real insane data",
			args: args{
				in: "37b92100",
			},
			want: "21b937",
		},
		{
			name: "test weird 0 pad",
			args: args{
				in: "00000000a0e83e01",
			},
			want: "013ee8a0",
		},
		{
			name: "test real insane data",
			args: args{
				in: "0264081501cb7195013461b7",
				//   013461B7 01CB7195 02640815
			},
			want: "013461b701cb719502640815",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeDataSane(tt.args.in); got != tt.want {
				t.Errorf("makeDataSane() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseInput(t *testing.T) {
	now := time.Now()
	type args struct {
		in       string
		timeSent time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []FridgeData
		wantErr bool
	}{
		{
			name: "test 1 node",
			args: args{
				in:       "b7613401",
				timeSent: now,
			},
			want: []FridgeData{
				FridgeData{
					FridgeID:     "2",
					Time:         now,
					Temperature:  2.1,
					Humidity:     1,
					DoorOpenings: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "test 3 nodes",
			args: args{
				in:       "0264081501cb7195013461b7",
				timeSent: now,
			},
			want: []FridgeData{
				FridgeData{
					FridgeID:     "2",
					Time:         now,
					Temperature:  2.1,
					Humidity:     1,
					DoorOpenings: 3,
				},
				FridgeData{
					FridgeID:     "3",
					Time:         now,
					Temperature:  1.1,
					Humidity:     1,
					DoorOpenings: 1,
				},
				FridgeData{
					FridgeID:     "4",
					Time:         now,
					Temperature:  1.1,
					Humidity:     1,
					DoorOpenings: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.args.in, tt.args.timeSent)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
