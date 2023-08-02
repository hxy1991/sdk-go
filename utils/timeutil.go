package utils

import "time"

func Now() time.Time {
	if IsPre() || IsProduction() {
		return time.Now()
	}

	return time.Now().Add(time.Duration(EnvVar.FakeTimeOffset) * time.Second)
}
