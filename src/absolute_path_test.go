package dotfiles

import (
	"path"
	"strings"
	"testing"
)

func TestAbsolutePathNew(t *testing.T) {
	a, err := NewAbsolutePath("/path/to/somewhere")
	if err != nil {
		t.Error(err)
	}

	if string(a) != "/path/to/somewhere" {
		t.Errorf("Invalid absolute path: %s", a)
	}

	if !a.Compare("/path/to/somewhere") {
		t.Errorf("Compare() must be used for comparing absolute path and string")
	}
}

func TestAbsolutePathHomeDirExpansion(t *testing.T) {
	a, err := NewAbsolutePath("~/foo")
	if err != nil {
		t.Error(err)
	}
	if !path.IsAbs(string(a)) {
		t.Errorf("'~' must be converted to full home directory path: %s", a)
	}
	if !strings.HasSuffix(string(a), "/foo") {
		t.Errorf("Invalid path conversion for home directory expansion: %s", a)
	}
}

func TestAbsolutePathEmpty(t *testing.T) {
	a, err := NewAbsolutePath("")
	if err != nil {
		t.Error(err)
	}
	if string(a) != "" {
		t.Errorf("Path converted from empty path must be empty: '%s'", a)
	}
	if !a.IsEmpty() {
		t.Errorf("IsEmpty() returns false when empty absolute path: '%s'", a)
	}
}

func TestAbsolutePathRelativePathError(t *testing.T) {
	_, err := NewAbsolutePath("relative-path")
	if err == nil {
		t.Errorf("NewAbsolutePath() must raise an error when relative path is given")
	}
}
