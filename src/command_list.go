package dotfiles

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func List(repo string) error {
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

	m, err := GetMappings(path.Join(repo, ".dotfiles"))
	if err != nil {
		return err
	}

	p, err := NewAbsolutePath(repo)
	if err != nil {
		return err
	}

	links, err := m.ActualLinks(p)
	if err != nil {
		return err
	}

	for source, dist := range links {
		fmt.Printf("'%s' -> '%s'\n", source, dist)
	}

	return nil
}
