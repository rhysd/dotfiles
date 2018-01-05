package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/dotfiles/src"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	cli = kingpin.New("dotfiles", "A dotfiles symlinks manager")

	clone       = cli.Command("clone", "Clone remote repository")
	clone_repo  = clone.Arg("repository", "Repository.  Format: 'user', 'user/repo-name', 'git@somewhere.com:repo.git, 'https://somewhere.com/repo.git'").Required().String()
	clone_path  = clone.Arg("path", "Path where repository cloned").String()
	clone_https = clone.Flag("https", "Use https:// instead of git@ protocol for `git clone`.").Short('h').Bool()

	link           = cli.Command("link", "Put symlinks to setup your configurations")
	link_dryrun    = link.Flag("dry", "Show what happens only").Bool()
	link_repo      = link.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()
	link_specified = link.Arg("files", "Files to link. If you specify no file, all will be linked.").Strings()
	// TODO link_no_default = link.Flag("no-default", "Link files specified by mappings.json and mappings_*.json")

	list      = cli.Command("list", "Show a list of symbolic link put by this command")
	list_repo = list.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()

	clean      = cli.Command("clean", "Remove all symbolic links put by this command")
	clean_repo = clean.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()

	update      = cli.Command("update", "Update your dotfiles repository")
	update_repo = update.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()

	version    = cli.Command("version", "Show version")
	updateSelf = cli.Command("selfupdate", "Show version")
)

func unimplemented(cmd string) {
	fmt.Fprintf(os.Stderr, "Command '%s' is not implemented yet!\n", cmd)
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		// Note: Exit code is detemined with looking http://tldp.org/LDP/abs/html/exitcodes.html
		os.Exit(113)
	} else {
		os.Exit(0)
	}
}

func selfUpdate() {
	v := semver.MustParse(dotfiles.Version())

	latest, err := selfupdate.UpdateSelf(v, "rhysd/dotfiles")
	if err != nil {
		exit(err)
	}

	if v.Equals(latest.Version) {
		fmt.Println("Current version", v, "is the latest")
	} else {
		fmt.Println("Successfully updated to version", v)
		fmt.Println("Release Note:\n", latest.ReleaseNotes)
	}
}

func main() {
	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case clone.FullCommand():
		exit(dotfiles.Clone(*clone_repo, *clone_path, *clone_https))
	case link.FullCommand():
		exit(dotfiles.Link(*link_repo, *link_specified, *link_dryrun))
	case list.FullCommand():
		exit(dotfiles.List(*list_repo))
	case clean.FullCommand():
		exit(dotfiles.Clean(*clean_repo))
	case update.FullCommand():
		exit(dotfiles.Update(*update_repo))
	case version.FullCommand():
		fmt.Println(dotfiles.Version())
	case updateSelf.FullCommand():
		selfUpdate()
	default:
		panic("Internal error: Unreachable! Please report this to https://github.com/rhysd/dotfiles/issues")
	}
}
