package dotfiles

import (
	"fmt"
)

func Clone(spec, specified string, https bool) error {
	repo, err := NewRepository(spec, specified, https)
	if err != nil {
		return err
	}

	err = repo.Clone()
	if err != nil {
		return err
	}

	s := "into"
	if repo.IncludesRepoDir {
		s = "as"
	}
	fmt.Printf("\nYour dotfiles was successfully cloned from '%s' %s '%s'\n", repo.Url, s, repo.Path.String())

	return nil
}
