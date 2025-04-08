package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// Level represents log severity
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
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
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Config holds logger configuration
type Config struct {
	// Level is the minimum log level to output
	Level Level
	// EnableJSON enables JSON formatted logs
	EnableJSON bool
	// EnableColor enables colored output for console logging
	EnableColor bool
	// EnableCaller adds caller information to logs
	EnableCaller bool
	// Output is where logs are written
	Output io.Writer
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Level:        InfoLevel,
		EnableJSON:   false,
		EnableColor:  true,
		EnableCaller: true,
		Output:       os.Stdout,
	}
}

// Logger handles application logging
type Logger struct {
	config Config
}

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	return &Logger{
		config: config,
	}
}

// Entry is a log entry with fields
type Entry struct {
	logger  *Logger
	level   Level
	message string
	fields  map[string]interface{}
	err     error
	caller  string
	time    time.Time
}

// Log creates a new entry with a message
func (l *Logger) Log(level Level, message string) *Entry {
	// Skip logging if level is too low
	if level < l.config.Level {
		return &Entry{logger: l}
	}

	e := &Entry{
		logger:  l,
		level:   level,
		message: message,
		fields:  make(map[string]interface{}),
		time:    time.Now(),
	}

	if l.config.EnableCaller {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			e.caller = fmt.Sprintf("%s:%d", file, line)
		}
	}

	return e
}

// Debug logs at debug level
func (l *Logger) Debug(message string) *Entry {
	return l.Log(DebugLevel, message)
}

// Info logs at info level
func (l *Logger) Info(message string) *Entry {
	return l.Log(InfoLevel, message)
}

// Warn logs at warn level
func (l *Logger) Warn(message string) *Entry {
	return l.Log(WarnLevel, message)
}

// Error logs at error level
func (l *Logger) Error(message string) *Entry {
	return l.Log(ErrorLevel, message)
}

// Fatal logs at fatal level and exits
func (l *Logger) Fatal(message string) *Entry {
	return l.Log(FatalLevel, message)
}

// WithField adds a field to the entry
func (e *Entry) WithField(key string, value interface{}) *Entry {
	// If entry is nil or below the configured log level, don't do anything
	if e.logger == nil {
		return e
	}

	e.fields[key] = value
	return e
}

// WithFields adds multiple fields to the entry
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	// If entry is nil or below the configured log level, don't do anything
	if e.logger == nil {
		return e
	}

	for k, v := range fields {
		e.fields[k] = v
	}
	return e
}

// WithError adds an error to the entry
func (e *Entry) WithError(err error) *Entry {
	// If entry is nil or below the configured log level, don't do anything
	if e.logger == nil {
		return e
	}

	e.err = err
	return e
}

// WithContext adds context values to the entry
func (e *Entry) WithContext(ctx context.Context) *Entry {
	// If entry is nil or below the configured log level, don't do anything
	if e.logger == nil {
		return e
	}

	// Add requestID if it exists
	if requestID, ok := ctx.Value("requestID").(string); ok {
		e.fields["requestID"] = requestID
	}

	return e
}

// Send writes the log entry
func (e *Entry) Send() {
	// If entry is nil or below the configured log level, don't do anything
	if e.logger == nil {
		return
	}

	if e.logger.config.EnableJSON {
		e.sendJSON()
	} else {
		e.sendText()
	}

	// If fatal, exit
	if e.level == FatalLevel {
		os.Exit(1)
	}
}

// sendJSON formats and sends a JSON log
func (e *Entry) sendJSON() {
	data := map[string]interface{}{
		"level":     e.level.String(),
		"message":   e.message,
		"timestamp": e.time.Format(time.RFC3339),
	}

	// Add fields
	for k, v := range e.fields {
		data[k] = v
	}

	// Add error if present
	if e.err != nil {
		data["error"] = e.err.Error()
	}

	// Add caller if enabled
	if e.logger.config.EnableCaller && e.caller != "" {
		data["caller"] = e.caller
	}

	// Encode as JSON
	encoder := json.NewEncoder(e.logger.config.Output)
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding log entry: %v\n", err)
	}
}

// sendText formats and sends a text log
func (e *Entry) sendText() {
	var levelColor, resetColor string
	if e.logger.config.EnableColor {
		resetColor = "\033[0m"
		switch e.level {
		case DebugLevel:
			levelColor = "\033[37m" // White
		case InfoLevel:
			levelColor = "\033[32m" // Green
		case WarnLevel:
			levelColor = "\033[33m" // Yellow
		case ErrorLevel:
			levelColor = "\033[31m" // Red
		case FatalLevel:
			levelColor = "\033[35m" // Purple
		}
	}

	// Format timestamp
	timestamp := e.time.Format("2006-01-02 15:04:05")

	// Build the log line
	fmt.Fprintf(
		e.logger.config.Output,
		"%s %s%s%s %s",
		timestamp,
		levelColor,
		e.level.String(),
		resetColor,
		e.message,
	)

	// Add fields if any
	if len(e.fields) > 0 {
		fmt.Fprint(e.logger.config.Output, " ")
		first := true
		for k, v := range e.fields {
			if first {
				first = false
			} else {
				fmt.Fprint(e.logger.config.Output, " ")
			}
			fmt.Fprintf(e.logger.config.Output, "%s=%v", k, v)
		}
	}

	// Add error if present
	if e.err != nil {
		fmt.Fprintf(e.logger.config.Output, " error=%q", e.err.Error())
	}

	// Add caller if enabled
	if e.logger.config.EnableCaller && e.caller != "" {
		fmt.Fprintf(e.logger.config.Output, " caller=%s", e.caller)
	}

	fmt.Fprintln(e.logger.config.Output)
}

// Global logger instance
var std = New(DefaultConfig())

// SetLevel sets the log level for the standard logger
func SetLevel(level Level) {
	std.config.Level = level
}

// SetOutput sets the output for the standard logger
func SetOutput(output io.Writer) {
	std.config.Output = output
}

// EnableJSON enables JSON formatting for the standard logger
func EnableJSON() {
	std.config.EnableJSON = true
}

// Debug logs at debug level
func Debug(message string) *Entry {
	return std.Debug(message)
}

// Info logs at info level
func Info(message string) *Entry {
	return std.Info(message)
}

// Warn logs at warn level
func Warn(message string) *Entry {
	return std.Warn(message)
}

// Error logs at error level
func Error(message string) *Entry {
	return std.Error(message)
}

// Fatal logs at fatal level and exits
func Fatal(message string) *Entry {
	return std.Fatal(message)
}
