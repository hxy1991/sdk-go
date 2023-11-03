package utils

import (
	"os"
	"strconv"
	"strings"
)

var EnvVar struct {
	FakeTimeOffset int
}

var serverEnv = os.Getenv("SERVER_ENV")

func IsProduction() bool {
	return strings.HasPrefix(strings.ToLower(serverEnv), "prod")
}

func IsPre() bool {
	return strings.HasPrefix(strings.ToLower(serverEnv), "pre")
}

func IsDevelopment() bool {
	return strings.HasPrefix(strings.ToLower(serverEnv), "dev")
}

func IsTest() bool {
	return strings.HasPrefix(strings.ToLower(serverEnv), "test")
}

func IsConsoleLog() bool {
	return strings.EqualFold(os.Getenv("CONSOLE_LOG"), "1")
}

func init() {
	if os.Getenv("FAKE_TIME_OFFSET") != "" {
		offset, err := strconv.Atoi(os.Getenv("FAKE_TIME_OFFSET"))
		if err != nil {
			offset = 0
		}
		EnvVar.FakeTimeOffset = offset
	}
}
