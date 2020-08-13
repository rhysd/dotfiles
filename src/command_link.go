package dotfiles

func Link(repoInput string, specified []string, dry bool) error {
	repo, err := absolutePathToRepo(repoInput)
	if err != nil {
		return err
	}

	m, err := GetMappings(repo.Join(".dotfiles"))
	if err != nil {
		return err
	}

	if len(specified) == 0 {
		err = m.CreateAllLinks(repo, dry)
		if e, ok := err.(*NothingLinkedError); ok {
			e.RepoPath = repo.String()
		}
		return err
	}

	return m.CreateSomeLinks(specified, repo, dry)
}
