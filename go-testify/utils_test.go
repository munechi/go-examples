package myutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHoge(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{2, 3, 5},
		{1, -1, 0},
		{10, -4, 6},
		{-4, 10, 6},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, Add(tt.a, tt.b))
	}
}
