package dotfiles

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

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
		r, err := absolutePathToRepo(c.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %s", c.input, err.Error())
			continue
		}
		if r.String() != c.expected {
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
		_, err := absolutePathToRepo(e)
		if err == nil {
			t.Errorf("'%s' is an invalid value for repository but no error occurred", e)
		}
	}

	for _, c := range []struct {
		env      string
		expected string
	}{
		{cwd, cwd},
		{"~", u.HomeDir},
		{".", abs(".")},
	} {
		os.Setenv("DOTFILES_REPO_PATH", c.env)
		r, err := absolutePathToRepo("")
		if err != nil {
			t.Errorf("Unexpected error for $DOEFILES_REPO_PATH '%s': %s", c.env, err.Error())
			continue
		}
		if r.String() != c.expected {
			t.Errorf("Expected '%s' as absolute path but actually '%s'", c.expected, r)
		}
	}
	os.Setenv("DOTFILES_REPO_PATH", "")
}
