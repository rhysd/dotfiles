package dotfiles

import (
	"os"
	"path"
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

	error_cases := []string{
		"LICENSE.txt",
		fullPath("LICENSE.txt"),
		"unknown",
		fullPath("unknown"),
	}

	for _, i := range error_cases {
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

	success_cases := map[string]string{
		"":                          getwd(),
		"_test_directory":           "_test_directory",
		fullPath("_test_directory"): fullPath("_test_directory"),
	}

	for input, expected := range success_cases {
		r, err := NewRepository("foo", input, false)
		if err != nil {
			t.Errorf("Unexpected error on specifying path: %s", err.Error())
		} else if r.ParentDir != expected {
			t.Errorf("Expected %s as the parent directory but actually %s", expected, r.ParentDir)
		}
	}
}

func TestNewRepositoryNormalizeRepoUrl(t *testing.T) {
	success_cases := map[string]string{
		"rhysd":                                 "git@github.com:rhysd/dotfiles.git",
		"rhysd/foobar":                          "git@github.com:rhysd/foobar.git",
		"https://github.com/rhysd/dogfiles.git": "https://github.com/rhysd/dogfiles.git",
		"git@bitbucket.com:rhysd/dotfiles.git":  "git@bitbucket.com:rhysd/dotfiles.git",
		"https://github.com/rhysd/dogfiles":     "https://github.com/rhysd/dogfiles.git",
		"git@bitbucket.com:rhysd/dotfiles":      "git@bitbucket.com:rhysd/dotfiles.git",
	}

	for input, expected := range success_cases {
		r, err := NewRepository(input, "", false)
		if err != nil {
			t.Errorf("Unexpected error for full path: %s: %s", input, err.Error())
		}
		if r.Url != expected {
			t.Errorf("Expected %s for input %s, but actually %s", expected, input, r.Url)
		}
	}
}

func TestNewRepositoryWithHttps(t *testing.T) {
	success_cases := map[string]string{
		"rhysd":                                 "https://github.com/rhysd/dotfiles.git",
		"rhysd/foobar":                          "https://github.com/rhysd/foobar.git",
		"https://github.com/rhysd/dogfiles.git": "https://github.com/rhysd/dogfiles.git",
		"git@bitbucket.com:rhysd/dotfiles.git":  "git@bitbucket.com:rhysd/dotfiles.git",
	}

	for input, expected := range success_cases {
		r, err := NewRepository(input, "", true)
		if err != nil {
			t.Errorf("Unexpected error for full path: %s: %s", input, err.Error())
		}
		if r.Url != expected {
			t.Errorf("Expected %s for input %s, but actually %s", expected, input, r.Url)
		}
	}
}

func TestNewRepositoryInvalidEmptySpec(t *testing.T) {
	_, err := NewRepository("", "", false)
	if err == nil {
		t.Errorf("Expected an error when empty spec was provided")
	}
}

func TestClone(t *testing.T) {
	if err := os.MkdirAll("_test_cloned", os.ModeDir|os.ModePerm); err != nil {
		panic(err.Error())
	}
	defer os.RemoveAll("_test_cloned")

	{
		r, _ := NewRepository("rhysd/vim-rustpeg", "", false)
		if err := r.Clone(); err != nil {
			t.Fatalf("Error on cloning repository %s to current directory: %s", r.Url, err.Error())
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

	{
		r, _ := NewRepository("git@bitbucket.org:rhysd/dotfiles", "_test_cloned", false)
		if err := r.Clone(); err != nil {
			t.Fatalf("Error on cloning repository %s to current directory: %s", r.Url, err.Error())
		}
		p := path.Join(getwd(), "_test_cloned", "dotfiles") // Just a test repository
		s, err := os.Stat(p)
		if err != nil {
			t.Fatalf("Cloned repository not found")
		}
		if !s.IsDir() {
			t.Fatalf("Cloned repository is not a directory")
		}
	}
}
