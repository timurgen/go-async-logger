package logmonkey

import "testing"

func TestGetLevelByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want LogLevel
	}{
		{name: "valid-1", args: args{"TRACE"}, want: TRACE},
		{name: "valid-2", args: args{"DebuG"}, want: DEBUG},
		{name: "valid-3", args: args{"infO"}, want: INFO},
		{name: "valid-4", args: args{"Warning"}, want: WARNING},
		{name: "valid-5", args: args{"error"}, want: ERROR},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLevelByName(tt.args.name); got != tt.want {
				t.Errorf("GetLevelByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
