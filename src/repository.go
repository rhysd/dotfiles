package dotfiles

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rhysd/abspath"
)

// Repository represents a repository on local filesystem
// TODO: Enable to specify branch name?
type Repository struct {
	URL             string
	Path            abspath.AbsPath
	IncludesRepoDir bool
}

func pathToCloneRepo(specified string) (abspath.AbsPath, bool, error) {
	if specified != "" {
		repo, err := abspath.ExpandFrom(specified)
		if err != nil {
			return abspath.AbsPath{}, false, err
		}
		if s, err := os.Stat(repo.String()); err != nil || !s.IsDir() {
			return abspath.AbsPath{}, false, fmt.Errorf("Specified path does not exist or is a file: '%s'", repo.String())
		}
		return repo, false, nil
	}

	if env := os.Getenv("DOTFILES_REPO_PATH"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return abspath.AbsPath{}, false, fmt.Errorf("Repository directory is specified as '%s' with $DOTFILES_REPO_PATH but it already exists", env)
		}
		repo, err := abspath.New(env)
		if err != nil {
			return abspath.AbsPath{}, false, err
		}
		return repo, true, nil
	}

	repo, err := abspath.ExpandFrom(".")
	if err != nil {
		return abspath.AbsPath{}, false, err
	}
	return repo, false, nil
}

func NewRepository(spec, specified string, https bool) (*Repository, error) {
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

	p, b, err := pathToCloneRepo(specified)
	if err != nil {
		return nil, err
	}
	return &Repository{spec, p, b}, nil
}

func (repo *Repository) Clone() error {
	args := []string{"clone", repo.URL}
	if repo.IncludesRepoDir {
		args = append(args, repo.Path.String())
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if cwd != repo.Path.String() {
			if err := os.Chdir(repo.Path.String()); err != nil {
				return err
			}
			defer os.Chdir(cwd)
		}
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
