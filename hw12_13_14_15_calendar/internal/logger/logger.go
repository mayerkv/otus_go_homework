package logger

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelDict = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

type Logger struct {
	lvl Level
	w   io.Writer
}

func New(lvl Level, w io.Writer) *Logger {
	return &Logger{lvl: lvl, w: w}
}

func (l *Logger) Debug(msg string) {
	l.Log(LevelDebug, msg)
}

func (l *Logger) Info(msg string) {
	l.Log(LevelInfo, msg)
}

func (l *Logger) Warn(msg string) {
	l.Log(LevelWarn, msg)
}

func (l *Logger) Error(msg string) {
	l.Log(LevelError, msg)
}

func (l *Logger) Log(level Level, msg string) {
	if l.lvl > level {
		return
	}

	prefix := "LOG"
	if p, ok := levelDict[level]; ok {
		prefix = p
	}

	fmt.Fprintf(l.w, "%s [%s] %s\n", time.Now().UTC().Format(time.RFC3339), prefix, msg)
}

func LevelFromString(l string) Level {
	for k, v := range levelDict {
		if strings.ToUpper(l) == v {
			return k
		}
	}
	return LevelInfo
}
