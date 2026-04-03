package mlbapi

import (
	"regexp"
	"testing"
)

func TestVersionFormat(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	if !pattern.MatchString(Version) {
		t.Fatalf("Version = %q, want semantic version tag format", Version)
	}
}
