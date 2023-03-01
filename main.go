package main

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	dotfiles "github.com/rhysd/dotfiles/src"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cli = kingpin.New("dotfiles", "A dotfiles symlinks manager")

	clone      = cli.Command("clone", "Clone remote repository")
	cloneRepo  = clone.Arg("repository", "Repository.  Format: 'user', 'user/repo-name', 'git@somewhere.com:repo.git, 'https://somewhere.com/repo.git'").Required().String()
	clonePath  = clone.Arg("path", "Path where repository cloned").String()
	cloneHTTPS = clone.Flag("https", "Use https:// instead of git@ protocol for `git clone`.").Short('h').Bool()

	link          = cli.Command("link", "Put symlinks to setup your configurations")
	linkDryRun    = link.Flag("dry", "Show what happens only").Bool()
	linkRepo      = link.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()
	linkSpecified = link.Arg("files", "Files to link. If you specify no file, all will be linked.").Strings()
	// TODO link_no_default = link.Flag("no-default", "Link files specified by mappings.json and mappings_*.json")

	list     = cli.Command("list", "Show a list of symbolic link put by this command")
	listRepo = list.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()

	clean     = cli.Command("clean", "Remove all symbolic links put by this command")
	cleanRepo = clean.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()

	update     = cli.Command("update", "Update your dotfiles repository")
	updateRepo = update.Arg("repo", "Path to your dotfiles repository.  If omitted, $DOTFILES_REPO_PATH is searched and fallback into the current directory.").String()

	version    = cli.Command("version", "Show version")
	updateSelf = cli.Command("selfupdate", "Update the executable binary by downloading the latest version from GitHub releases page.")
)

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
		exit(dotfiles.Clone(*cloneRepo, *clonePath, *cloneHTTPS))
	case link.FullCommand():
		exit(dotfiles.Link(*linkRepo, *linkSpecified, *linkDryRun))
	case list.FullCommand():
		exit(dotfiles.List(*listRepo))
	case clean.FullCommand():
		exit(dotfiles.Clean(*cleanRepo))
	case update.FullCommand():
		exit(dotfiles.Update(*updateRepo))
	case version.FullCommand():
		fmt.Println(dotfiles.Version())
	case updateSelf.FullCommand():
		selfUpdate()
	default:
		panic("Internal error: Unreachable! Please report this to https://github.com/rhysd/dotfiles/issues")
	}
}
