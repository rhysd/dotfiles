package dotfiles

import (
	"fmt"
)

func Clone(spec, path string) error {
	repo, err := NewRepository(spec, path)
	if err != nil {
		return err
	}

	err = repo.Clone()
	if err != nil {
		return err
	}

	fmt.Printf("\nYour dotfiles was successfully cloned from '%s' into '%s'\n", repo.Url, repo.ParentDir)

	return nil
}
