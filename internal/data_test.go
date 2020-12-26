package internal

import (
	"reflect"
	"testing"
)

func TestNewChannelData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want *ChannelData
	}{
		{
			name: "Pos1",
			args: args{raw: []byte{0x00, 0x25, 0x46}},
			want: &ChannelData{
				Temperature: 3.7,
				Humidity:    70,
			},
		},
		{
			name: "Pos2",
			args: args{raw: []byte{0x01, 0x64, 0x63}},
			want: &ChannelData{
				Temperature: 35.6,
				Humidity:    99,
			},
		},
		{
			name: "Neg1",
			args: args{raw: []byte{0xFF, 0xFB, 0x09}},
			want: &ChannelData{
				Temperature: -0.5,
				Humidity:    9,
			},
		},
		{
			name: "Neg2",
			args: args{raw: []byte{0xFF, 0x6A, 0x27}},
			want: &ChannelData{
				Temperature: -15,
				Humidity:    39,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelData(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelData_RawBytes(t *testing.T) {
	tests := []struct {
		name string
		have *ChannelData
		want []byte
	}{
		{
			name: "Pos1",
			have: &ChannelData{
				Temperature: 3.7,
				Humidity:    70,
			},
			want: []byte{0x00, 0x25, 0x46},
		},
		{
			name: "Pos2",
			have: &ChannelData{
				Temperature: 35.6,
				Humidity:    99,
			},
			want: []byte{0x01, 0x64, 0x63},
		},
		{
			name: "Neg1",
			have: &ChannelData{
				Temperature: -0.5,
				Humidity:    9,
			},
			want: []byte{0xFF, 0xFB, 0x09},
		},
		{
			name: "Neg2",
			have: &ChannelData{
				Temperature: -15,
				Humidity:    39,
			},
			want: []byte{0xFF, 0x6A, 0x27},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.have
			if got := c.RawBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewChannelsData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want []*ChannelData
	}{
		{
			name: "2Channels",
			args: args{
				raw: []byte{0x00, 0x25, 0x46, 0x01, 0x64, 0x63},
			},
			want: []*ChannelData{
				{
					Number:      1,
					Temperature: 3.7,
					Humidity:    70,
				},
				{
					Number:      2,
					Temperature: 35.6,
					Humidity:    99,
				},
			},
		},
		{
			name: "4Channels",
			args: args{
				raw: []byte{0x00, 0x25, 0x46, 0x01, 0x64, 0x63, 0xFF, 0xFB, 0x09, 0xFF, 0x6A, 0x27},
			},
			want: []*ChannelData{
				{
					Number:      1,
					Temperature: 3.7,
					Humidity:    70,
				},
				{
					Number:      2,
					Temperature: 35.6,
					Humidity:    99,
				},
				{
					Number:      3,
					Temperature: -0.5,
					Humidity:    9,
				},
				{
					Number:      4,
					Temperature: -15,
					Humidity:    39,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelsData(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelsData() = %v, want %v", got, tt.want)
			}
		})
	}
}
