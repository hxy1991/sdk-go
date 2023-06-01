package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsProduction(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsProduction(), "IsProduction()")
		})
	}
}

func TestIsTest(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsTest(), "IsTest()")
		})
	}
}
