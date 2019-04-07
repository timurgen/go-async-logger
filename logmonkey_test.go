package logmonkey

import (
	"reflect"
	"testing"
	"time"
)

func TestGetLevelByName(t *testing.T) {
	defer FlushAllLoggers()
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

func TestLogLevel_String(t *testing.T) {
	defer FlushAllLoggers()
	tests := []struct {
		name string
		l    LogLevel
		want string
	}{
		{name: "TRACE", l: 0, want: "TRACE"},
		{name: "DEBUG", l: 1, want: "DEBUG"},
		{name: "INFO", l: 2, want: "INFO"},
		{name: "WARNING", l: 3, want: "WARNING"},
		{name: "ERROR", l: 4, want: "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultLogFormatter_FormatMessage(t *testing.T) {
	defer FlushAllLoggers()
	type args struct {
		message string
		name    string
		level   LogLevel
		ts      time.Time
	}
	tests := []struct {
		name string
		lf   *DefaultLogFormatter
		args args
		want string
	}{
		{
			name: "defult",
			lf:   &DefaultLogFormatter{Format: "%s - [%s] %s %s"},
			args: args{
				message: "test message",
				name:    "main",
				level:   1,
				ts:      time.Date(2019, 4, 1, 18, 0, 0, 0, time.UTC),
			},
			want: "2019-04-01T18:00:00.000000000 - [main] DEBUG test message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lf.FormatMessage(tt.args.message, tt.args.name, tt.args.level, tt.args.ts); got != tt.want {
				t.Errorf("DefaultLogFormatter.FormatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsoleLogAppender_ConsumeMessage(t *testing.T) {
	defer FlushAllLoggers()
	type args struct {
		str string
	}
	tests := []struct {
		name string
		la   *ConsoleLogAppender
		args args
	}{
		{
			name: "default",
			la:   &ConsoleLogAppender{},
			args: args{
				str: "test message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.la.ConsumeMessage(tt.args.str)
		})
	}
}

func TestLogger_SetAppender(t *testing.T) {
	defer FlushAllLoggers()
	type args struct {
		l LogAppender
	}
	tests := []struct {
		name string
		log  *Logger
		args args
	}{
		{
			name: "default",
			log:  GetLogger("test"),
			args: args{
				l: &ConsoleLogAppender{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.log.SetAppender(tt.args.l)
		})
	}
}

func TestLogger_SetFormatter(t *testing.T) {
	defer FlushAllLoggers()
	type args struct {
		f LogFormatter
	}
	tests := []struct {
		name string
		log  *Logger
		args args
	}{
		{
			name: "name",
			log:  GetLogger("test"),
			args: args{
				f: &DefaultLogFormatter{Format: "%s - [%s] %s %s"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.log.SetFormatter(tt.args.f)
		})
	}
}

func TestLogger_SetLevel(t *testing.T) {
	defer FlushAllLoggers()
	type args struct {
		level LogLevel
	}
	tests := []struct {
		name string
		log  *Logger
		args args
	}{
		{
			name: "test",
			log:  GetLogger("test"),
			args: args{
				level: DEBUG,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.log.SetLevel(tt.args.level)
		})
	}
}

func TestLogger_GetLevel(t *testing.T) {
	defer FlushAllLoggers()
	tests := []struct {
		name string
		log  *Logger
		want LogLevel
	}{
		{
			name: "test default log level",
			log:  GetLogger("test"),
			want: INFO,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.log.GetLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logger.GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
