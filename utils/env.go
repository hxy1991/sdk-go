package utils

import (
	"os"
	"strings"
)

func IsProduction() bool {
	return strings.EqualFold(os.Getenv("SERVER_ENV"), "prod")
}

func IsPre() bool {
	return strings.EqualFold(os.Getenv("SERVER_ENV"), "pre")
}

func IsDevelopment() bool {
	return strings.EqualFold(os.Getenv("SERVER_ENV"), "dev")
}

func IsTest() bool {
	return strings.EqualFold(os.Getenv("SERVER_ENV"), "test")
}
