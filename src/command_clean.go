package dotfiles

import "path/filepath"

func Clean(repo_input string) error {
	repo, err := AbsolutePathToRepo(repo_input)
	if err != nil {
		return err
	}

	m, err := GetMappings(filepath.Join(string(repo), ".dotfiles"))
	if err != nil {
		return err
	}

	return m.UnlinkAll(repo)
}
