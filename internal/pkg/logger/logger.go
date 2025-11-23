package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	mu     sync.Mutex
	output io.Writer
	level  Level
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Caller    string                 `json:"caller,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(os.Stdout, DebugLevel)
}

func New(output io.Writer, level Level) *Logger {
	return &Logger{
		output: output,
		level:  level,
	}
}

func SetLevel(levelStr string) {
	var level Level
	switch levelStr {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	default:
		level = DebugLevel
	}
	defaultLogger.level = level
}

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   msg,
		Caller:    caller,
		Fields:    fields,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	fmt.Fprintln(l.output, string(data))
}

func Debug(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields)
	defaultLogger.log(DebugLevel, msg, f)
}

func Info(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields)
	defaultLogger.log(InfoLevel, msg, f)
}

func Warn(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields)
	defaultLogger.log(WarnLevel, msg, f)
}

func Error(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields)
	defaultLogger.log(ErrorLevel, msg, f)
}

func mergeFields(fields []map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	for _, f := range fields {
		for k, v := range f {
			result[k] = v
		}
	}
	return result
}

func WithFields(fields map[string]interface{}) map[string]interface{} {
	return fields
}
