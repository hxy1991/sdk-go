package opensearchv2

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New(context.Background(), "test")
	if err != nil {
		t.Error(err)
	}
}
