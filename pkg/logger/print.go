package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var (
	logger   *slog.Logger
	logLevel = slog.LevelInfo
	initOnce = false
)

// ANSI color codes
var levelColors = map[slog.Level]string{
	slog.LevelDebug: "\033[36m", // Cyan
	slog.LevelInfo:  "\033[32m", // Green
	slog.LevelWarn:  "\033[33m", // Yellow
	slog.LevelError: "\033[31m", // Red
}

const colorReset = "\033[0m"

// CustomHandler implements slog.Handler with custom formatting
type CustomHandler struct {
	w     io.Writer
	level slog.Level
}

func NewCustomHandler(w io.Writer, level slog.Level) *CustomHandler {
	return &CustomHandler{w: w, level: level}
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomHandler) Handle(_ context.Context, record slog.Record) error {
	// Format time as 20250730-0815.xxx
	timeStr := record.Time.Format("20060102-1504.000")

	// Get level string
	levelStr := record.Level.String()

	// Get color for level
	color := levelColors[record.Level]

	// Format: "20250730-0815.xxx INFO message"
	formatted := fmt.Sprintf("%s %s%s%s %s\n",
		timeStr,
		color,
		levelStr,
		colorReset,
		record.Message)

	_, err := h.w.Write([]byte(formatted))
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we'll ignore attributes in this custom format
	// You can extend this if you need to handle attributes
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll ignore groups in this custom format
	// You can extend this if you need to handle groups
	return h
}

func Init() {
	h := NewCustomHandler(os.Stderr, logLevel)
	logger = slog.New(h)
	initOnce = true
}

// SetLogLevel sets the logger's level. Call before Init().
func SetLogLevel(level string) {
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	if initOnce {
		// Re-initialize logger with new level
		Init()
	}
}

func ensureInit() {
	if !initOnce {
		Init()
	}
}

func makeLogFunc(level slog.Level, logFunc func(string, ...any)) func(string, ...any) {
	return func(format string, a ...any) {
		ensureInit()
		if logger.Enabled(context.TODO(), level) {
			logFunc(format, a...)
		}
	}
}

var (
	LogDebug = makeLogFunc(slog.LevelDebug, func(format string, a ...any) {
		logger.Debug(fmt.Sprintf(format, a...))
	})
	LogInfo = makeLogFunc(slog.LevelInfo, func(format string, a ...any) {
		logger.Info(fmt.Sprintf(format, a...))
	})
	LogWarn = makeLogFunc(slog.LevelWarn, func(format string, a ...any) {
		logger.Warn(fmt.Sprintf(format, a...))
	})
	LogError = makeLogFunc(slog.LevelError, func(format string, a ...any) {
		logger.Error(fmt.Sprintf(format, a...))
	})
)
