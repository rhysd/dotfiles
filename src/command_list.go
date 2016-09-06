package dotfiles

import (
	"fmt"
	"path/filepath"
)

func List(specified string) error {
	repo, err := AbsolutePathToRepo(specified)
	if err != nil {
		return err
	}

	m, err := GetMappings(filepath.Join(string(repo), ".dotfiles"))
	if err != nil {
		return err
	}

	links, err := m.ActualLinks(repo)
	if err != nil {
		return err
	}

	for source, dist := range links {
		fmt.Printf("'%s' -> '%s'\n", source, dist)
	}

	if len(links) == 0 {
		fmt.Printf("No link was found (dotfiles: %s)\n", repo)
	}

	return nil
}
