package log

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Debug(_ context.Context, message string, args ...interface{}) {
	write(log.Debug(), message, args)
}

func Info(_ context.Context, message string, args ...interface{}) {
	write(log.Info(), message, args)
}

func Warn(_ context.Context, message string, args ...interface{}) {
	write(log.Warn(), message, args)
}

func Fatal(_ context.Context, message string, args ...interface{}) {
	write(log.Fatal(), message, args)
}

func Error(_ context.Context, message string, args ...interface{}) {
	write(log.Error(), message, args)
}

func write(evt *zerolog.Event, message string, args ...interface{}) {
	for i, arg := range args {
		evt = evt.Interface(numbers[i], arg)
	}

	evt.Msg(message)
}

func Debugf(_ context.Context, message string, args ...interface{}) {
	writef(log.Debug(), message, args)
}

func Infof(_ context.Context, message string, args ...interface{}) {
	writef(log.Info(), message, args)
}

func Warnf(_ context.Context, message string, args ...interface{}) {
	writef(log.Warn(), message, args)
}

func Fatalf(_ context.Context, message string, args ...interface{}) {
	writef(log.Fatal(), message, args)
}

func Errorf(_ context.Context, message string, args ...interface{}) {
	writef(log.Error(), message, args)
}

func writef(evt *zerolog.Event, message string, args ...interface{}) {
	evt.Msgf(message, args...)
}

var numbers = map[int]string{
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "10",
}
