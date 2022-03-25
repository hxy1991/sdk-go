package utils

import "os"

func IsProduction() bool {
	return os.Getenv("SERVER_ENV") == "prod"
}

func IsPre() bool {
	return os.Getenv("SERVER_ENV") == "pre"
}

func IsDevelopment() bool {
	return os.Getenv("SERVER_ENV") == "dev"
}

func IsTest() bool {
	return os.Getenv("SERVER_ENV") == "test"
}
