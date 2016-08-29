package dotfiles

import (
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	r := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	if !r.MatchString(v) {
		t.Fatalf("Invalid versioning: %s", v)
	}
}
