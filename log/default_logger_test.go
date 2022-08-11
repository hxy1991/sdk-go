package log

import "testing"

func TestInfo(t *testing.T) {
	Info("1", "2", "3")
}

func TestWithMap(t *testing.T) {
	map1 := map[string]interface{}{
		"url":      "1",
		"userId":   1,
		"actionId": 2,
		"helperId": 3,
		"now":      "4",
	}

	WithMap(map1).Info()
}
