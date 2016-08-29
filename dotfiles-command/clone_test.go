package dotfiles

import (
	"os"
	"testing"
)

func TestCloneCommand(t *testing.T) {
	if err := Clone("rhysd/vim-rustpeg", ""); err != nil {
		t.Fatalf("Unexpected error on cloning: %s", err.Error())
	}
	defer os.RemoveAll("vim-rustpeg")
	s, err := os.Stat("vim-rustpeg")
	if err != nil {
		t.Fatalf("Cloned repository not found")
	}
	if !s.IsDir() {
		t.Fatalf("Cloned repository is not a directory")
	}
}
