package main

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
)

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
