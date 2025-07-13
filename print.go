package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	logger *log.Logger
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

func init() {
	// Create custom logger with 3-digit microsecond precision
	customWriter := &customLogger{writer: os.Stderr}
	logger = log.New(customWriter, "", 0) // No flags since we handle timestamp ourselves
}

// Only prints when Verbose is true
func VPrintln(a ...interface{}) {
	if Verbose {
		logger.Println(a...)
	}
}
func VPrintf(format string, a ...interface{}) {
	if Verbose {
		logger.Printf(format, a...)
	}
}
func VPrint(a ...interface{}) {
	if Verbose {
		logger.Print(a...)
	}
}
