package dotfiles

import (
	"fmt"
	"os"
	"path"
)

func Link(repo string, specified []string, dry bool) error {
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
