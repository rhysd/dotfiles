package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
)

func Clean(repo string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if repo == "" {
		repo = cwd
	} else {
		if !filepath.IsAbs(repo) {
			repo = filepath.Join(cwd, repo)
		}
		s, err := os.Stat(repo)
		if err != nil {
			return err
		}
		if !s.IsDir() {
			return fmt.Errorf("'%s' is not a directory. Please specify your dotfiles directory.", repo)
		}
	}

	p, err := NewAbsolutePath(repo)
	if err != nil {
		return err
	}

	m, err := GetMappings(filepath.Join(string(p), ".dotfiles"))
	if err != nil {
		return err
	}

	return m.UnlinkAll(p)
}
