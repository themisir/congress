package logger

import (
	"fmt"
	"io"
	"os"
)

type LogLevel byte

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

var levelNames = map[LogLevel]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
}

type Logger struct {
	name     string
	minLevel LogLevel
}

var Default = &Logger{"", Debug}

func New(name string, level LogLevel) *Logger {
	return &Logger{name, level}
}

func (l *Logger) Child(name string) *Logger {
	return &Logger{name, l.minLevel}
}

func (l *Logger) Log(level LogLevel, format string, v ...interface{}) {
	if level >= l.minLevel {
		var output io.Writer
		if level < Error {
			output = os.Stdout
		} else {
			output = os.Stderr
		}
		if l.name != "" {
			fmt.Fprintf(output, "[%s] %s: ", levelNames[level], l.name)
		} else {
			fmt.Fprintf(output, "[%s] ", levelNames[level])
		}
		fmt.Fprintf(output, format+"\n", v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.Log(Debug, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Log(Info, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.Log(Warn, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Log(Error, format, v...)
}
