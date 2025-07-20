package utils

import (
	"os"
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Convert absolute path to relative path from home directory
func ToHomeRelativePath(absPath string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return absPath // Return original path if we can't get home directory
	}

	// Check if the path starts with the home directory
	if strings.HasPrefix(absPath, homeDir) {
		// Replace home directory with ~
		relativePath := strings.Replace(absPath, homeDir, "~", 1)
		return relativePath
	}

	return absPath // Return original path if it's not under home directory
}
