package main

import (
	"fmt"
	"github.com/rhysd/dotfiles/dotfiles-command"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app = kingpin.New("dotfiles", "A dotfiles manager")

	clone      = app.Command("clone", "Clone remote repository")
	clone_repo = clone.Arg("repository", "Repository.  Format: 'user', 'user/repo-name', 'git@somewhere.com:repo.git, 'https://somewhere.com/repo.git'").Required().String()
	clone_path = clone.Arg("path", "Path where repository cloned").String()

	link        = app.Command("link", "Put symlinks to setup your configurations")
	link_dryrun = link.Flag("dry", "Show what happens only").Bool()

	list = app.Command("list", "Show a list of symbolic link put by this command")

	clean = app.Command("clean", "Remove all symbolic links put by this command")

	update = app.Command("update", "Update your dotfiles repository")

	version = app.Command("version", "Show version")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case clone.FullCommand():
		fmt.Printf("clone: %s %v '%s'\n", clone.FullCommand(), *clone_repo, *clone_path)
	case link.FullCommand():
		fmt.Printf("link: %s %v\n", link.FullCommand(), *link_dryrun)
	case list.FullCommand():
		fmt.Printf("list: %s\n", list.FullCommand())
	case clean.FullCommand():
		fmt.Printf("clean: %s\n", clean.FullCommand())
	case update.FullCommand():
		fmt.Printf("update: %s\n", update.FullCommand())
	case version.FullCommand():
		fmt.Println(dotfiles.Version())
	default:
		fmt.Printf("unknown\n")
	}
}
