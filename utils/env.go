package utils

import (
	"os"
	"strings"
)

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
