package log

import (
	"github.com/fatih/color"
	"os"
)

type Logger struct {
	info  *color.Color
	debug *color.Color
	error *color.Color
	fatal *color.Color
}

func New() *Logger {
	return &Logger{
		info:  color.New(color.FgCyan),
		debug: color.New(color.FgHiBlack),
		error: color.New(color.FgRed),
		fatal: color.New(color.FgHiRed),
	}
}

func (l *Logger) Debugln(args ...interface{}) {
	l.debug.Println(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.debug.Printf(format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.info.Println(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.error.Println(args...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.error.Printf(format, args...)
	os.Exit(1)
}
