package dotfiles

import (
	"os"
	"path"
	"testing"
)

func TestLinkAll(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dist_conf := path.Join(cwd, "_dist.conf")
	dir := path.Join(cwd, ".dotfiles")
	if err := os.MkdirAll(dir, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	f, err := os.OpenFile(path.Join(dir, "mappings.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.RemoveAll(dir)
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(`
	{
		"_source.conf": "` + dist_conf + `"
	}
	`)
	if err != nil {
		panic(err)
	}

	source := path.Join(cwd, "_source.conf")
	g, err := os.OpenFile(source, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		g.Close()
		os.Remove(source)
	}()
	_, err = g.WriteString("this file is for test")
	if err != nil {
		panic(err)
	}

	if err := Link("", nil, false); err != nil {
		t.Error(err)
	}
	defer os.Remove("_dist.conf")
}

func TestLinkSome(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dist_conf := path.Join(cwd, "_dist.conf")
	dir := path.Join(cwd, ".dotfiles")
	if err := os.MkdirAll(dir, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	f, err := os.OpenFile(path.Join(dir, "mappings.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.RemoveAll(dir)
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(`
	{
		"_source.conf": "` + dist_conf + `",
		"_tmp.conf": "/path/to/somewhere"
	}
	`)
	if err != nil {
		panic(err)
	}

	source := path.Join(cwd, "_source.conf")
	g, err := os.OpenFile(source, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		g.Close()
		os.Remove(source)
	}()
	_, err = g.WriteString("this file is for test")
	if err != nil {
		panic(err)
	}

	if err := Link("", []string{"_source.conf"}, false); err != nil {
		t.Error(err)
	}
	defer os.Remove("_dist.conf")
}

func TestLinkConfigDirDoesNotExist(t *testing.T) {
	if err := Link("", nil, false); err != nil {
		if _, ok := err.(*NothingLinkedError); !ok {
			t.Errorf("Non-existtence of .dotfiles directory does not cause an error: %s", err.Error())
		}
	}
}

func TestLinkSpecifiedRepoDoesNotExist(t *testing.T) {
	if err := Link("unknown_directory", nil, false); err == nil {
		t.Errorf("Should make an error for unknown dotfiles repository")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	p := path.Join(cwd, "_dummy_file")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
		os.Remove(p)
	}()

	_, err = f.WriteString("dummy file")
	if err != nil {
		panic(err)
	}

	if err := Link("_dummy_file", nil, false); err == nil {
		t.Errorf("Should make an error when repository is actually a file")
	}
}
