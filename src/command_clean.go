package dotfiles

func Clean(repo_input string) error {
	repo, err := absolutePathToRepo(repo_input)
	if err != nil {
		return err
	}

	m, err := GetMappings(repo.Join(".dotfiles"))
	if err != nil {
		return err
	}

	return m.UnlinkAll(repo)
}
