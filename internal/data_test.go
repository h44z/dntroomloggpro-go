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

func TestNewCalibrationData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want *CalibrationData
	}{
		{
			name: "Pos1",
			args: args{raw: []byte{0x00, 0x00, 0x00}},
			want: &CalibrationData{
				Temperature: 0,
				Humidity:    0,
			},
		},
		{
			name: "Pos2",
			args: args{raw: []byte{0x01, 0x64, 0x63}},
			want: &CalibrationData{
				Temperature: 35.6,
				Humidity:    99,
			},
		},
		{
			name: "Neg1",
			args: args{raw: []byte{0xFF, 0xFB, 0x09}},
			want: &CalibrationData{
				Temperature: -0.5,
				Humidity:    9,
			},
		},
		{
			name: "Neg2",
			args: args{raw: []byte{0xFF, 0x6A, 0x27}},
			want: &CalibrationData{
				Temperature: -15,
				Humidity:    39,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCalibrationData(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCalibrationData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCalibrationsData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want []*CalibrationData
	}{
		{
			name: "2Channels",
			args: args{
				raw: []byte{0x00, 0x25, 0x46, 0x01, 0x64, 0x63},
			},
			want: []*CalibrationData{
				{
					Channel:     1,
					Temperature: 3.7,
					Humidity:    70,
				},
				{
					Channel:     2,
					Temperature: 35.6,
					Humidity:    99,
				},
			},
		},
		{
			name: "4Channels",
			args: args{
				raw: []byte{0x00, 0x25, 0x46, 0x00, 0x00, 0x00, 0xFF, 0xFB, 0x09, 0xFF, 0x6A, 0x27},
			},
			want: []*CalibrationData{
				{
					Channel:     1,
					Temperature: 3.7,
					Humidity:    70,
				},
				{
					Channel:     2,
					Temperature: 0,
					Humidity:    0,
				},
				{
					Channel:     3,
					Temperature: -0.5,
					Humidity:    9,
				},
				{
					Channel:     4,
					Temperature: -15,
					Humidity:    39,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCalibrationsData(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCalibrationsData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalibrationData_RawBytes(t *testing.T) {
	tests := []struct {
		name string
		have *CalibrationData
		want []byte
	}{
		{
			name: "Pos1",
			have: &CalibrationData{
				Temperature: 3.7,
				Humidity:    70,
			},
			want: []byte{0x00, 0x25, 0x46},
		},
		{
			name: "Pos2",
			have: &CalibrationData{
				Temperature: 35.6,
				Humidity:    99,
			},
			want: []byte{0x01, 0x64, 0x63},
		},
		{
			name: "Neg1",
			have: &CalibrationData{
				Temperature: -0.5,
				Humidity:    9,
			},
			want: []byte{0xFF, 0xFB, 0x09},
		},
		{
			name: "Neg2",
			have: &CalibrationData{
				Temperature: -15,
				Humidity:    39,
			},
			want: []byte{0xFF, 0x6A, 0x27},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.have
			if got := d.RawBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIntervalData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want IntervalData
	}{
		{
			name: "Zero",
			args: args{raw: []byte{0x00}},
			want: IntervalData(0),
		},
		{
			name: "Min",
			args: args{raw: []byte{0x01}},
			want: IntervalData(1),
		},
		{
			name: "Misc1",
			args: args{raw: []byte{0x21}},
			want: IntervalData(33),
		},
		{
			name: "Misc2",
			args: args{raw: []byte{0x63}},
			want: IntervalData(99),
		},
		{
			name: "Max",
			args: args{raw: []byte{0xff}},
			want: IntervalData(255),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIntervalData(tt.args.raw); got != tt.want {
				t.Errorf("NewIntervalData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntervalData_RawBytes(t *testing.T) {
	tests := []struct {
		name string
		d    IntervalData
		want []byte
	}{
		{
			name: "Zero",
			want: []byte{0x00},
			d:    IntervalData(0),
		},
		{
			name: "Min",
			want: []byte{0x01},
			d:    IntervalData(1),
		},
		{
			name: "Misc1",
			want: []byte{0x21},
			d:    IntervalData(33),
		},
		{
			name: "Misc2",
			want: []byte{0x63},
			d:    IntervalData(99),
		},
		{
			name: "Max",
			want: []byte{0xff},
			d:    IntervalData(255),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.RawBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSettingsAreaData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want *SettingsAreaData
	}{
		{
			name: "Ch1Temp",
			args: args{
				raw: []byte{0x00, 0x00, 0x00, 0x01},
			},
			want: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: true,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
		},
		{
			name: "Ch2TempCh2Dew",
			args: args{
				raw: []byte{0x00, 0x00, 0x00, 0x18},
			},
			want: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: false,
					1: true,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: true,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
		},
		{
			name: "Ch3TempCh4TempCh5Temp",
			args: args{
				raw: []byte{0x00, 0x00, 0x12, 0x40},
			},
			want: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: false,
					1: false,
					2: true,
					3: true,
					4: true,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
		},
		{
			name: "Ch3DewCh4DewCh5Dew",
			args: args{
				raw: []byte{0x00, 0x00, 0x24, 0x80},
			},
			want: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: false,
					2: true,
					3: true,
					4: true,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSettingsAreaData(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettingsAreaData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSettingsAreaData_RawBytes(t *testing.T) {
	tests := []struct {
		name string
		have *SettingsAreaData
		want []byte
	}{
		{
			name: "Ch1Temp",
			have: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: true,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
			want: []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "Ch2TempCh2Dew",
			have: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: false,
					1: true,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: true,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
			want: []byte{0x00, 0x00, 0x00, 0x18},
		},
		{
			name: "Ch3TempCh4TempCh5Temp",
			have: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: false,
					1: false,
					2: true,
					3: true,
					4: true,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
			want: []byte{0x00, 0x00, 0x12, 0x40},
		},
		{
			name: "Ch3DewCh4DewCh5Dew",
			have: &SettingsAreaData{
				Temperature: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
				DewPoint: map[uint8]bool{
					0: false,
					1: false,
					2: true,
					3: true,
					4: true,
					5: false,
					6: false,
					7: false,
				},
				HeatIndex: map[uint8]bool{
					0: false,
					1: false,
					2: false,
					3: false,
					4: false,
					5: false,
					6: false,
					7: false,
				},
			},
			want: []byte{0x00, 0x00, 0x24, 0x80},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.have
			if got := d.RawBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSettingsData(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want *SettingsData
	}{
		{
			name: "Settings1",
			args: args{
				raw: []byte{0x01, 0x48, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x00},
			},
			want: &SettingsData{
				GraphType:     GraphTypeHumidity,
				GraphInterval: GraphInterval72h,
				TimeFormat:    TimeFormatEurope,
				DateFormat:    DateFormatDDMMYYYY,
				DST:           DSTOff,
				TimeZone:      2,
				Units:         UnitCelsius,
				Areas: [5]*SettingsAreaData{
					{
						Area: 1,
						Temperature: map[uint8]bool{
							0: true,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 2,
						Temperature: map[uint8]bool{
							0: false,
							1: true,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 3,
						Temperature: map[uint8]bool{
							0: false,
							1: false,
							2: true,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 4,
						Temperature: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: true,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 5,
						Temperature: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: true,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSettingsData(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettingsData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSettingsData_RawBytes(t *testing.T) {
	tests := []struct {
		name string
		have *SettingsData
		want []byte
	}{
		{
			name: "Settings1",
			have: &SettingsData{
				GraphType:     GraphTypeHumidity,
				GraphInterval: GraphInterval72h,
				TimeFormat:    TimeFormatEurope,
				DateFormat:    DateFormatDDMMYYYY,
				DST:           DSTOff,
				TimeZone:      2,
				Units:         UnitCelsius,
				Areas: [5]*SettingsAreaData{
					{
						Area: 1,
						Temperature: map[uint8]bool{
							0: true,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 2,
						Temperature: map[uint8]bool{
							0: false,
							1: true,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 3,
						Temperature: map[uint8]bool{
							0: false,
							1: false,
							2: true,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 4,
						Temperature: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: true,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
					{
						Area: 5,
						Temperature: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: true,
							5: false,
							6: false,
							7: false,
						},
						DewPoint: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
						HeatIndex: map[uint8]bool{
							0: false,
							1: false,
							2: false,
							3: false,
							4: false,
							5: false,
							6: false,
							7: false,
						},
					},
				},
			},
			want: []byte{0x01, 0x48, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.have
			if got := d.RawBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
