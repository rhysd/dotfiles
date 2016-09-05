package dotfiles

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
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

func TestAbsolutePathToRepo(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	abs := func(p string) string {
		a, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}
		return a
	}

	for _, c := range []struct {
		input    string
		expected string
	}{
		{".", abs(".")},
		{"", cwd},
		{"../src", cwd},
		{"~", u.HomeDir},
		{cwd, cwd},
	} {
		r, err := AbsolutePathToRepo(c.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %s", c.input, err.Error())
			continue
		}
		if string(r) != c.expected {
			t.Errorf("Expected '%s' as absolute path but actually '%s'", c.expected, r)
		}
	}

	f, err := os.OpenFile("existing_file", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	f.Close()
	defer os.Remove("existing_file")

	for _, e := range []string{
		"unknown_dir",
		"./existing_file",
	} {
		_, err := AbsolutePathToRepo(e)
		if err == nil {
			t.Errorf("'%s' is an invalid value for repository but no error occurred", e)
		}
	}
}
