/*
	Copyright 2021 Misir Jafarov

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

			http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
		var line string
		if level < Error {
			output = os.Stdout
		} else {
			output = os.Stderr
		}
		if l.name != "" {
			line += fmt.Sprintf("[%s] %s: ", levelNames[level], l.name)
		} else {
			line += fmt.Sprintf("[%s] ", levelNames[level])
		}
		fmt.Fprintf(output, line+format+"\n", v...)
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
