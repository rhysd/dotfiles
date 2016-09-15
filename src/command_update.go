package dotfiles

import (
	"os"
	"os/exec"
)

func Update(repo_input string) error {
	repo, err := absolutePathToRepo(repo_input)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if repo.String() != cwd {
		if err := os.Chdir(repo.String()); err != nil {
			return err
		}
		defer os.Chdir(cwd)
	}

	cmd := exec.Command("git", "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
