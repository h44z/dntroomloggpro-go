package internal

import (
	"reflect"
	"testing"
)

func TestGetMessagePayload(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Valid1",
			args: args{
				raw: []byte{0x7b, 0x01, 0x02, 0x03, 0x40, 0x7d, 0x99, 0xff},
			},
			want:    []byte{0x01, 0x02, 0x03},
			wantErr: false,
		},
		{
			name: "Valid2",
			args: args{
				raw: []byte{0x7b, 0x01, 0x02, 0x03, 0x40, 0x7d},
			},
			want:    []byte{0x01, 0x02, 0x03},
			wantErr: false,
		},
		{
			name: "Valid3",
			args: args{
				raw: []byte{0x7b, 0x40, 0x7d},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "Invalid1",
			args: args{
				raw: []byte{0x7b, 0x00, 0x7d},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid2",
			args: args{
				raw: []byte{0x7b},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid3",
			args: args{
				raw: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMessagePayload(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessagePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMessagePayload() got = %v, want %v", got, tt.want)
			}
		})
	}
}
