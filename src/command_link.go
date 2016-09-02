package dotfiles

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func Link(repo string, specified []string, dry bool) error {
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

	if specified == nil || len(specified) == 0 {
		err = m.CreateAllLinks(dry)
		if e, ok := err.(*NothingLinkedError); ok {
			e.Repo = repo
		}
		return err
	} else {
		return m.CreateSomeLinks(specified, dry)
	}
}
