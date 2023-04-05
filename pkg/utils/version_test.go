package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGetCRDAVersion(t *testing.T) {
	assert.Equal(t, GetCRDAVersion(), version)
}

func TestBuildVersion(t *testing.T) {
	// example expected output: v0.0.0-abcd linux/amd64  BuildDate: Thu, 01 Jan 1970 02:00:00 IST  Vendor: Local Build
	targetRgx := fmt.Sprintf(
		"^v%s-%s [a-z]+/[a-z0-9]+  BuildDate: [A-Za-z]{3}, [0-9]{2} [A-Za-z]{3} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Z]{3}  Vendor: %s$",
		version,
		commitHash,
		vendorInfo,
	)
	assert.Regexp(t, regexp.MustCompile(targetRgx), BuildVersion())
}
