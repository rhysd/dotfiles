package dotfiles

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type AbsolutePath string

func NewAbsolutePath(s string) (AbsolutePath, error) {
	if s == "" {
		return AbsolutePath(""), nil
	}

	if s[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return AbsolutePath(""), err
		}
		return AbsolutePath(filepath.Join(u.HomeDir, s[1:])), nil
	}

	if !filepath.IsAbs(s) {
		return "", fmt.Errorf("Not an absolute path: '%s'", s)
	}
	return AbsolutePath(s), nil
}

func AbsolutePathToRepo(repo string) (AbsolutePath, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return AbsolutePath(""), err
	}

	if repo == "" {
		repo = cwd
	} else if repo[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return AbsolutePath(""), err
		}
		repo = filepath.Join(u.HomeDir, repo[1:])
	} else if !filepath.IsAbs(repo) {
		repo = filepath.Join(cwd, repo)
	}

	s, err := os.Stat(repo)
	if err != nil {
		return AbsolutePath(""), err
	}

	if !s.IsDir() {
		return AbsolutePath(""), fmt.Errorf("'%s' is not a directory. Please specify your dotfiles directory.", repo)
	}

	return AbsolutePath(repo), nil
}

func (a AbsolutePath) Compare(s string) bool {
	// Note: Should we consider '~' in s?
	return string(a) == s
}

func (a AbsolutePath) IsEmpty() bool {
	return len(a) == 0
}
