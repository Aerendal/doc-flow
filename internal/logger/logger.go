package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

func ParseFormat(s string) Format {
	if s == "json" {
		return FormatJSON
	}
	return FormatText
}

type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	level  Level
	format Format
	fields map[string]string
}

type jsonEntry struct {
	Time    string            `json:"time"`
	Level   string            `json:"level"`
	Message string            `json:"msg"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func New(out io.Writer, level Level, format Format) *Logger {
	return &Logger{
		out:    out,
		level:  level,
		format: format,
		fields: make(map[string]string),
	}
}

func NewDefault() *Logger {
	return New(os.Stderr, LevelInfo, FormatText)
}

func (l *Logger) WithField(key, value string) *Logger {
	newFields := make(map[string]string, len(l.fields)+1)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &Logger{
		out:    l.out,
		level:  l.level,
		format: l.format,
		fields: newFields,
	}
}

func (l *Logger) WithDocID(docID string) *Logger {
	return l.WithField("doc_id", docID)
}

func (l *Logger) WithPath(path string) *Logger {
	return l.WithField("path", path)
}

func (l *Logger) log(level Level, msg string, args ...any) {
	if level < l.level {
		return
	}

	formatted := fmt.Sprintf(msg, args...)

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format("2006-01-02T15:04:05")

	if l.format == FormatJSON {
		entry := jsonEntry{
			Time:    now,
			Level:   level.String(),
			Message: formatted,
		}
		if len(l.fields) > 0 {
			entry.Fields = l.fields
		}
		data, _ := json.Marshal(entry)
		fmt.Fprintln(l.out, string(data))
		return
	}

	fieldsStr := ""
	for k, v := range l.fields {
		fieldsStr += fmt.Sprintf(" %s=%s", k, v)
	}
	fmt.Fprintf(l.out, "%s [%s] %s%s\n", now, level.String(), formatted, fieldsStr)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args...)
}
