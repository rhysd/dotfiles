package dotfiles

import (
	"os"
	"path"
	"testing"
)

func TestCleanAll(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	repos := []string{"", cwd}

	for _, repo := range repos {
		f, err := os.OpenFile("_source.conf", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString("this file is for test")
		if err != nil {
			panic(err)
		}
		f.Close()
		defer os.Remove("_source.conf")

		linked := path.Join(cwd, "_linked.conf")
		if err := os.Symlink(path.Join(cwd, "_source.conf"), linked); err != nil {
			panic(err)
		}
		defer os.Remove(linked)

		config := path.Join(cwd, ".dotfiles")
		if err := os.MkdirAll(config, os.ModePerm|os.ModeDir); err != nil {
			panic(err)
		}
		defer os.RemoveAll(config)

		f, err = os.OpenFile(path.Join(config, "mappings.json"), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(`
		{
			"_source.conf": "` + linked + `"
		}
		`)
		if err != nil {
			panic(err)
		}
		f.Close()

		err = Clean(repo)
		if err != nil {
			t.Error(err)
		}

		if _, err := os.Lstat(linked); err == nil {
			t.Errorf("Unlinked symlink must be removed")
		}
	}
}

func TestCleanAllInvalidRepo(t *testing.T) {
	if err := Clean("unknown_dir"); err == nil {
		t.Errorf("Non-existing repository directory must raise an error")
	}

	f, err := os.OpenFile("file_as_repository", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	f.Close()
	defer os.Remove("file_as_repository")

	if err := Clean("file_as_repository"); err == nil {
		t.Errorf("Should raise an error when directory is actually a file")
	}
}
