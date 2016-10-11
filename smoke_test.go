package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestSmoke(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	repo := filepath.Join(cwd, "dogfiles")
	if err := os.Setenv("DOTFILES_REPO_PATH", repo); err != nil {
		panic(err)
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := user.HomeDir
	defer os.RemoveAll(repo)

	{
		cmd := exec.Command("./dotfiles", "clone", "--https", "rhysd/dogfiles")
		if out, err := cmd.Output(); err != nil {
			t.Fatalf("Error on 'clone': %s: %s", err.Error(), out)
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
		buf, err := exec.Command("../dotfiles", "link", "--dry").Output()
		if err != nil {
			t.Fatal(err)
		}
		s := string(buf[:])

		expected := fmt.Sprintf("Link:  'zshrc' -> '%s'", filepath.Join(home, ".zshrc"))
		if !strings.Contains(s, expected) {
			t.Errorf("'%s' must be included in 'link' output on --dry: %s", expected, s)
		}

		expected = fmt.Sprintf("Link:  'cabal_config' -> '%s'", filepath.Join(home, ".cabal", "config"))
		if !strings.Contains(s, expected) {
			t.Errorf("'%s' must be included in 'link' output on --dry: %s", expected, s)
		}
	}

	{
		if out, err := exec.Command("../dotfiles", "link").Output(); err != nil {
			t.Fatalf("Error on 'link': %s: %s", err.Error(), out)
		}
		for _, l := range []struct {
			from string
			to   string
		}{
			{"npmrc", filepath.Join(home, ".npmrc")},
			{"nvimrc", filepath.Join(home, ".config", "nvim", "init.vim")},
			{"tmux.conf", filepath.Join(home, ".tmux.conf")},
			{"peco", filepath.Join(home, ".config", "peco")},
			{"vimrc", filepath.Join(home, ".vimrc")},
			{"gvimrc", filepath.Join(home, ".gvimrc")},
			{"nyaovim", filepath.Join(home, ".config", "nyaovim")},
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

			dist, err := os.Readlink(l.to)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.HasSuffix(dist, l.from) {
				t.Errorf("Unexpected link from '%s' -> to '%s'. '%s' is expected as destination", l.to, dist, l.from)
			}
		}
	}

	{
		buf, err := exec.Command("../dotfiles", "list").Output()
		if err != nil {
			t.Fatal(err)
		}
		s := string(buf[:])
		p, err := filepath.Abs("zshrc")
		if err != nil {
			panic(err)
		}
		if !strings.Contains(s, p) {
			t.Errorf("'list' must contain the path to zshrc in dotfiles")
		}
		p, err = filepath.Abs("peco")
		if err != nil {
			panic(err)
		}
		if !strings.Contains(s, p) {
			t.Errorf("'list' must contain the path to peco directory in dotfiles")
		}
		if !strings.Contains(s, filepath.Join(home, ".zshrc")) {
			t.Errorf("'list' must contain the path to ~/.zshrc in dotfiles")
		}
		if !strings.Contains(s, filepath.Join(home, ".config", "nvim", "init.vim")) {
			t.Errorf("'list' must contain the path to init.vim")
		}
	}

	{
		if out, err := exec.Command("../dotfiles", "update").Output(); err != nil {
			t.Fatalf("Error on 'update': %s: %s", err.Error(), out)
		}
	}

	{
		if out, err := exec.Command("../dotfiles", "clean").Output(); err != nil {
			t.Fatalf("Error on 'clean': %s: %s", err.Error(), out)
		}
		for _, l := range []struct {
			from string
			to   string
		}{
			{"npmrc", filepath.Join(home, ".npmrc")},
			{"nvimrc", filepath.Join(home, ".config", "nvim", "init.vim")},
			{"tmux.conf", filepath.Join(home, ".tmux.conf")},
			{"peco", filepath.Join(home, ".config", "peco")},
			{"vimrc", filepath.Join(home, ".vimrc")},
			{"gvimrc", filepath.Join(home, ".gvimrc")},
			{"nyaovim", filepath.Join(home, ".config", "nyaovim")},
		} {
			_, err := os.Lstat(l.to)
			if err == nil {
				t.Errorf("Symbolic link '%s' must be removed", l.to)
				continue
			}
		}
	}
}
