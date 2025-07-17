package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logger     *log.Logger
	loggerOnce sync.Once
)

// Custom logger that writes with 3-digit microsecond precision
type customLogger struct {
	writer io.Writer
}

func (cl *customLogger) Write(p []byte) (n int, err error) {
	// Get current time with microsecond precision
	now := time.Now()
	// Format time with 3 digits of microseconds (divide by 1000 to get milliseconds, then format)
	timestamp := fmt.Sprintf("%s.%03d", now.Format("2006/01/02 15:04:05"), now.Nanosecond()/1000000)

	// Write timestamp + space + original message
	formatted := fmt.Sprintf("%s %s", timestamp, string(p))
	return cl.writer.Write([]byte(formatted))
}

func getLogger() *log.Logger {
	loggerOnce.Do(func() {
		customWriter := &customLogger{writer: os.Stderr}
		logger = log.New(customWriter, "", 0) // No flags since we handle timestamp ourselves
	})
	return logger
}

// Only prints when Verbose is true
func VPrintln(a ...interface{}) {
	if Verbose {
		getLogger().Println(a...)
	}
}
func VPrintf(format string, a ...interface{}) {
	if Verbose {
		getLogger().Printf(format, a...)
	}
}
func VPrint(a ...interface{}) {
	if Verbose {
		getLogger().Print(a...)
	}
}

// Always-on info log (use for important info, not just verbose)
func LogInfo(format string, a ...interface{}) {
	getLogger().Printf(format, a...)
}
