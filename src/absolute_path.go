package dotfiles

import (
	"fmt"
	"github.com/rhysd/abspath"
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

func AbsolutePathToRepo(repo string) (abspath.AbsPath, error) {
	if repo == "" {
		repo = os.Getenv("DOTFILES_REPO_PATH")
	}

	if repo == "" {
		repo = "."
		fmt.Fprintln(os.Stderr, "No repository was specified nor $DOTFILES_REPO_PATH was not set. Assuming current repository is a dotfiles repository.\n")
	}

	p, err := abspath.ExpandFrom(repo)
	if err != nil {
		return abspath.AbsPath{}, err
	}

	if !p.IsDir() {
		return abspath.AbsPath{}, fmt.Errorf("'%s' is not a directory. Please specify your dotfiles directory.", p.String())
	}

	return p, nil
}

func (a AbsolutePath) Compare(s string) bool {
	// Note: Should we consider '~' in s?
	return string(a) == s
}

func (a AbsolutePath) IsEmpty() bool {
	return len(a) == 0
}
