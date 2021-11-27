package log

import (
	"github.com/fatih/color"
	"os"
)

type Logger struct {
	info  *color.Color
	debug *color.Color
	warn  *color.Color
	fatal *color.Color
}

func New() *Logger {
	return &Logger{
		info:  color.New(color.FgCyan),
		debug: color.New(color.FgHiBlack),
		warn:  color.New(color.FgRed),
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

func (l *Logger) Warnln(args ...interface{}) {
	l.warn.Println(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.warn.Printf(format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.fatal.Println(args...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.fatal.Printf(format, args...)
	os.Exit(1)
}
