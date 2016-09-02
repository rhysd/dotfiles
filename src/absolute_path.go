package dotfiles

import (
	"fmt"
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

func (a AbsolutePath) Compare(s string) bool {
	// Note: Should we consider '~' in s?
	return string(a) == s
}

func (a AbsolutePath) IsEmpty() bool {
	return len(a) == 0
}
