package dotfiles

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TODO: Enable to specify branch name?
type Repository struct {
	Url       string
	ParentDir string
}

func NewRepository(spec, path string, https bool) (*Repository, error) {
	if path == "" {
		var err error
		if path, err = os.Getwd(); err != nil {
			return nil, err
		}
	} else {
		s, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if !s.IsDir() {
			return nil, fmt.Errorf("'%s' is not a directory", path)
		}
	}
	if spec == "" {
		return nil, fmt.Errorf("Remote path to clone must not be empty")
	}
	if strings.HasPrefix(spec, "https://") {
		if !strings.HasSuffix(spec, ".git") {
			spec = spec + ".git"
		}
	} else if strings.HasPrefix(spec, "git@") {
		if !strings.HasSuffix(spec, ".git") {
			spec = spec + ".git"
		}
	} else if strings.ContainsRune(spec, '/') {
		if https {
			spec = fmt.Sprintf("https://github.com/%s.git", spec)
		} else {
			spec = fmt.Sprintf("git@github.com:%s.git", spec)
		}
	} else {
		if https {
			spec = fmt.Sprintf("https://github.com/%s/dotfiles.git", spec)
		} else {
			spec = fmt.Sprintf("git@github.com:%s/dotfiles.git", spec)
		}
	}
	return &Repository{spec, path}, nil
}

func (repo *Repository) Clone() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if repo.ParentDir != cwd {
		if err := os.Chdir(repo.ParentDir); err != nil {
			return err
		}
		defer os.Chdir(cwd)
	}

	cmd := exec.Command("git", "clone", repo.Url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
