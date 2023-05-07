package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchSnykRegex(t *testing.T) {
	t.Run("empty string as token should return false", func(t *testing.T) {
		assert.False(t, MatchSnykRegex(""))
	})

	t.Run("legit token should return true", func(t *testing.T) {
		assert.True(t, MatchSnykRegex("a012345B-123C-432d-B123-80123456789B"))
	})
}
