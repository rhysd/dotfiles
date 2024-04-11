package dotfiles

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func getwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	return cwd
}

func fullPath(entry string) string {
	return path.Join(getwd(), entry)
}

func TestNewRepositoryInvalidPath(t *testing.T) {
	var err error

	errorCases := []string{
		"LICENSE.txt",
		fullPath("LICENSE.txt"),
		"unknown",
		fullPath("unknown"),
	}

	for _, i := range errorCases {
		_, err = NewRepository("foo", i, false)
		if err == nil {
			t.Errorf("Error was expected with invalid input to <path>:%s", i)
		}
	}
}

func TestNewRepositoryValidPath(t *testing.T) {
	if err := os.MkdirAll("_test_directory", os.ModeDir|os.ModePerm); err != nil {
		panic(err.Error())
	}
	defer os.Remove("_test_directory")

	successCases := map[string]string{
		"":                          getwd(),
		"_test_directory":           "_test_directory",
		fullPath("_test_directory"): fullPath("_test_directory"),
	}

	for input, expected := range successCases {
		r, err := NewRepository("foo", input, false)
		if err != nil {
			t.Errorf("Unexpected error on specifying path: %s", err.Error())
		} else if !strings.HasSuffix(r.Path.String(), expected) {
			t.Errorf("Expected %s as the parent directory but actually %s", expected, r.Path)
		}
	}
}

func TestNewRepositoryNormalizeRepoUrl(t *testing.T) {
	successCases := map[string]string{
		"rhysd":                                 "git@github.com:rhysd/dotfiles.git",
		"rhysd/foobar":                          "git@github.com:rhysd/foobar.git",
		"https://github.com/rhysd/dogfiles.git": "https://github.com/rhysd/dogfiles.git",
		"https://github.com/rhysd/dogfiles":     "https://github.com/rhysd/dogfiles.git",
	}

	for input, expected := range successCases {
		r, err := NewRepository(input, "", false)
		if err != nil {
			t.Errorf("Unexpected error for full path: %s: %s", input, err.Error())
		}
		if r.URL != expected {
			t.Errorf("Expected %s for input %s, but actually %s", expected, input, r.URL)
		}
	}
}

func TestNewRepositoryWithHttps(t *testing.T) {
	successCases := map[string]string{
		"rhysd":                                 "https://github.com/rhysd/dotfiles.git",
		"rhysd/foobar":                          "https://github.com/rhysd/foobar.git",
		"https://github.com/rhysd/dogfiles.git": "https://github.com/rhysd/dogfiles.git",
	}

	for input, expected := range successCases {
		r, err := NewRepository(input, "", true)
		if err != nil {
			t.Errorf("Unexpected error for full path: %s: %s", input, err.Error())
		}
		if r.URL != expected {
			t.Errorf("Expected %s for input %s, but actually %s", expected, input, r.URL)
		}
	}
}

func TestNewRepositoryInvalidEmptySpec(t *testing.T) {
	_, err := NewRepository("", "", false)
	if err == nil {
		t.Errorf("Expected an error when empty spec was provided")
	}
}

func TestNewRepositoryWithEnv(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	repo := filepath.Join(cwd, "_test_dotfiles")
	os.Setenv("DOTFILES_REPO_PATH", repo)
	defer os.Setenv("DOTFILES_REPO_PATH", "")

	r, err := NewRepository("rhysd/dogfiles", "", true)
	if err != nil {
		t.Fatal(err)
	}

	if !r.IncludesRepoDir {
		t.Errorf("dotfiles always includes its repository name")
	}

	if r.Path.String() != repo {
		t.Errorf("Repository must be installed at %s but actually done at %s", repo, r.Path.String())
	}
}

func TestNewRepositoryWithInvalidEnv(t *testing.T) {
	defer os.Setenv("DOTFILES_REPO_PATH", "")

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	os.Setenv("DOTFILES_REPO_PATH", cwd)
	if _, err := NewRepository("rhysd/dogfiles", "", true); err == nil {
		t.Fatalf("It must raise an error when repository already exists")
	}

	os.Setenv("DOTFILES_REPO_PATH", "this_is_relative_path")
	if _, err := NewRepository("rhysd/dogfiles", "", true); err == nil {
		t.Fatalf("It must raise an error when relative path is specified in env")
	}
}

func TestCloneDummyCommand(t *testing.T) {
	saved := os.Getenv("DOTFILES_GIT_COMMAND")
	defer os.Setenv("DOTFILES_GIT_COMMAND", saved)
	if err := os.Setenv("DOTFILES_GIT_COMMAND", "true"); err != nil {
		panic(err)
	}

	r, _ := NewRepository("rhysd/vim-rustpeg", "", false)
	if err := r.Clone(); err != nil {
		t.Fatalf("Error on cloning repository %s to current directory: %s", r.URL, err)
	}
}

func TestCloneError(t *testing.T) {
	r, _ := NewRepository("rhysd/repository-does-not-exist", "", false)
	if err := r.Clone(); err == nil {
		t.Fatalf("Error did not occur")
	}
}

func TestCloneRepo(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skip test for cloning repository on GitHub Actions")
	}

	{
		r, _ := NewRepository("rhysd/vim-rustpeg", "", false)
		if err := r.Clone(); err != nil {
			t.Fatalf("Error on cloning repository %s to current directory: %s", r.URL, err.Error())
		}
		defer os.RemoveAll("vim-rustpeg")
		p := path.Join(getwd(), "vim-rustpeg") // Just a test repository
		s, err := os.Stat(p)
		if err != nil {
			t.Fatalf("Cloned repository not found")
		}
		if !s.IsDir() {
			t.Fatalf("Cloned repository is not a directory")
		}
	}

	if err := os.MkdirAll("_test_cloned", os.ModeDir|os.ModePerm); err != nil {
		panic(err.Error())
	}
	defer os.RemoveAll("_test_cloned")
}

func TestCloneWithEnv(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	repo := filepath.Join(cwd, "_test_dotfiles")
	os.Setenv("DOTFILES_REPO_PATH", repo)
	defer os.Setenv("DOTFILES_REPO_PATH", "")

	r, err := NewRepository("rhysd/dogfiles", "", true)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Clone(); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("_test_dotfiles")

	s, err := os.Stat(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !s.IsDir() {
		t.Fatalf("Cloned repository must be a directory: '%s'", repo)
	}
}
