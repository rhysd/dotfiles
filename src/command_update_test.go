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
