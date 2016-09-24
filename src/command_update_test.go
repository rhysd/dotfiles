package dotfiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateErrorCase(t *testing.T) {
	if err := Update("unknown_repo"); err == nil {
		t.Fatalf("It should raise an error when unknown repository specified")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := Update(filepath.Base(cwd)); err == nil {
		t.Fatalf("If it is not a Git repository, it should raise an error")
	}
}

func TestUpdateOk(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := Update(".."); err != nil {
		t.Fatal(err)
	}
	if c, _ := os.Getwd(); c != cwd {
		t.Fatalf("Current working directory is wrong. '%s' should be '%s'", c, cwd)
	}
}
