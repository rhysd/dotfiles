package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.SetEnv("DOTFILES_REPO_PATH", cwd); err != nil {
		panic(err)
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := user.HomeDir

	{
		cmd := exec.Command("./dotfiles", "clone", "--https", "rhysd/dogfiles")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		s, err := os.Stat("dogfiles")
		if err != nil {
			t.Fatal(err)
		}
		if !s.IsDir() {
			t.Fatalf("Cloned repository is not a directory")
		}

		os.Chdir("dogfiles")
	}

	{
		buf, err := exec.Command("../dotfiles", "list").Output()
		if err != nil {
			t.Fatal(err)
		}
		s := buf.String()
		if !strings.Contains(filepath.Abs("zshrc")) {
			t.Errorf("'list' must contain the path to zshrc in dotfiles")
		}
		if !strings.Contains(filepath.Abs("peco")) {
			t.Errorf("'list' must contain the path to peco directory in dotfiles")
		}
		if !strings.Contains(filepath.Join(home, ".zshrc")) {
			t.Errorf("'list' must contain the path to ~/.zshrc in dotfiles")
		}
		if !strings.Contains(filepath.join(home, ".config", "nvim", "init.vim")) {
			t.Errorf("'list' must contain the path to init.vim")
		}
	}

	{
		buf, err := exec.Command("../dotfiles", "link", "--dry").Output()
		if err != nil {
			t.Fatal(err)
		}
		s := buf.String()

		expected := fmt.Sprintf("Link:  'zshrc' -> '%s'", filepath.Join(home, ".zshrc"))
		if !strings.Contains(expected) {
			t.Errorf("'%s' must be included in 'link' output on --dry", expected)
		}

		expected = fmt.Sprintf("Link:  'cabal_config' -> '%s'", filepath.Join(home, ".cabal", "config"))
		if !strings.Contains(expected) {
			t.Errorf("'%s' must be included in 'link' output on --dry", expected)
		}
	}

	{
		if err := exec.Command("../dotfiles", "link").Run(); err != nil {
			t.Fatalf(err)
		}
		for _, l := range []struct {
			from string
			to   string
		}{
			"npmrc":     filepath.Join(home, ".npmrc"),
			"nvimrc":    filepath.Join(home, ".config", "nvim", "init.vim"),
			"tmux.conf": filepath.Join(home, ".tmux.conf"),
			"peco":      filepath.Join(home, ".config", "peco"),
			"vimrc":     filepath.Join(home, ".vimrc"),
			"gvimrc":    filepath.Join(home, ".gvimrc"),
			"nyaovim":   filepath.Join(home, ".config", "nyaovim"),
		} {
			s, err := os.Lstat(l.to)
			if err != nil {
				t.Error(err)
				continue
			}

			if s.Mode()&os.ModeSymlink != os.ModeSymlink {
				t.Errorf("'%s' is not a symbolic link", l.to)
				continue
			}

			dist, err := os.Readlink(source)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.HasSuffix(dist, l.from) {
				t.Errorf("Unexpected link from '%s' -> to '%s'. '%s' is expected as destination", l.to, dist, l.from)
			}
		}
	}

	{
		if err := exec.Command("../dotfiles", "update").Run(); err != nil {
			t.Fatalf(err)
		}
	}

	{
		if err := exec.Command("../dotfiles", "clean").Run(); err != nil {
			t.Fatalf(err)
		}
		for _, l := range []struct {
			from string
			to   string
		}{
			"npmrc":     filepath.Join(home, ".npmrc"),
			"nvimrc":    filepath.Join(home, ".config", "nvim", "init.vim"),
			"tmux.conf": filepath.Join(home, ".tmux.conf"),
			"peco":      filepath.Join(home, ".config", "peco"),
			"vimrc":     filepath.Join(home, ".vimrc"),
			"gvimrc":    filepath.Join(home, ".gvimrc"),
			"nyaovim":   filepath.Join(home, ".config", "nyaovim"),
		} {
			s, err := os.Lstat(l.to)
			if err == nil {
				t.Error("Symbolic link '%s' must be removed", l.to)
				continue
			}
		}
	}
}
