package mlbapi

import (
	"regexp"
	"testing"
)

func TestVersionFormat(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`^(dev|v\d+\.\d+\.\d+)$`)
	if !pattern.MatchString(Version) {
		t.Fatalf("Version = %q, want dev or semantic version tag format", Version)
	}
}
