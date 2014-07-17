package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockNow(t *testing.T) {
	fixed := ""
	FixedTime("2014-07-01 00:00:00 +0900", func() {
		fixed = Now().Format("2006-01-02 15:04:05")
		assert.Equal(t, "2014-07-01 00:00:00", fixed)
	})
	assert.NotEqual(t, fixed, Now().Format("2006-01-02 15:04:05"))
}
