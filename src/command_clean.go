package dotfiles

import (
	"fmt"
	"os"
	"path"
)

func Clean(repo string) error {
	if repo == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		repo = cwd
	} else {
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
