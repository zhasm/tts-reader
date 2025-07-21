package logger

import "testing"

func TestInitAndLogInfo(t *testing.T) {
	Init()
	LogInfo("test log: %d", 123)
}

func TestVerboseLogging(t *testing.T) {
	Init()
	Verbose = true
	VPrintln("test VPrintln")
	VPrintf("test VPrintf: %d", 42)
	VPrint("test VPrint")
}
