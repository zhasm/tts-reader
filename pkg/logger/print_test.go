package logger

import (
	"testing"
)

func TestInitAndLogInfo(t *testing.T) {
	Init()
	LogInfo("test log: %d", 123)
}
