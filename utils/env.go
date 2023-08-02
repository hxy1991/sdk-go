package utils

import (
	"os"
	"strconv"
	"strings"
)

var EnvVar struct {
	FakeTimeOffset int
}

func IsProduction() bool {
	return strings.HasPrefix(strings.ToLower(os.Getenv("SERVER_ENV")), "prod")
}

func IsPre() bool {
	return strings.HasPrefix(strings.ToLower(os.Getenv("SERVER_ENV")), "pre")
}

func IsDevelopment() bool {
	return strings.HasPrefix(strings.ToLower(os.Getenv("SERVER_ENV")), "dev")
}

func IsTest() bool {
	return strings.HasPrefix(strings.ToLower(os.Getenv("SERVER_ENV")), "test")
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
