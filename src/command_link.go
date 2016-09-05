package dotfiles

import "path/filepath"

func Link(repo_input string, specified []string, dry bool) error {
	repo, err := AbsolutePathToRepo(repo_input)
	if err != nil {
		return err
	}

	m, err := GetMappings(filepath.Join(string(repo), ".dotfiles"))
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
