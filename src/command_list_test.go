package dotfiles

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListEmptyList(t *testing.T) {
	stdout_saved := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = stdout_saved
	}()

	if err := List("."); err != nil {
		t.Fatal(err)
	}
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	s := buf.String()
	if len(s) > 0 {
		t.Errorf("When no valid mapping exists, it should output nothing, but output '%s'", s)
	}
}

func TestExistingMapping(t *testing.T) {
	stdout_saved := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = stdout_saved
	}()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dist_conf := filepath.Join(cwd, "_dist.conf")
	dir := filepath.Join(cwd, ".dotfiles")
	if err := os.MkdirAll(dir, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	f, err := os.OpenFile(filepath.Join(dir, "mappings.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(`
	{
		"_source.conf": "` + dist_conf + `"
	}
	`)
	if err != nil {
		panic(err)
	}
	f.Close()

	source := filepath.Join(dir, "_source.conf")
	g, err := os.OpenFile(source, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	_, err = g.WriteString("this file is for test")
	if err != nil {
		panic(err)
	}
	g.Close()

	if err := os.Symlink(source, dist_conf); err != nil {
		panic(err)
	}
	defer os.Remove(dist_conf)

	if err := List(""); err != nil {
		t.Fatal(err)
	}
	w.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	if !strings.Contains(s, source) {
		t.Errorf("Output must contains source file path: '%s'", s)
	}
	if !strings.Contains(s, dist_conf) {
		t.Errorf("Output must contains dist symlink path: '%s'", s)
	}
}

func TestListInvalidInput(t *testing.T) {
	if err := List("/path/to/unknown_dir"); err == nil {
		t.Errorf("Unknown repository must raise an error")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir := filepath.Join(cwd, ".dotfiles")
	if err := os.MkdirAll(dir, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	f, err := os.OpenFile(filepath.Join(dir, "mappings.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(`
	{
		"broken_json":
	`)
	if err != nil {
		panic(err)
	}
	f.Close()

	if err := List("."); err == nil {
		t.Errorf("Broken JSON should raise an error on getting mappings.")
	}
}
