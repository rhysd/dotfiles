`dotfiles` Command
==================
[![CI](https://github.com/rhysd/dotfiles/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/rhysd/dotfiles/actions?query=workflow%3ACI)
[![Coverage](https://codecov.io/gh/rhysd/dotfiles/branch/master/graph/badge.svg)](https://codecov.io/gh/rhysd/dotfiles)

This repository provides `dotfiles` command to manage your [dotfiles](http://dotfiles.github.io/).  It manages your dotfiles repository and symbolic links to use the configurations.

This command has below goals:

- **One binary executable**: If you want to set configuration files in a remote server, all you have to do is sending a binary to the remote.
- **Do one thing and do it well**: This command manages only a dotfiles repository.  Does not handle any other dependencies.  If you want full-setup including dependencies, you should use more suitable tool such as [Ansible](https://www.ansible.com/).  And then use `dotfiles` command from it.
- **Less dependency**: Only depends on `git` command.
- **Sensible defaults**: Many sensible default symbolic link mappings are pre-defined.  You need not to specify the mappings for almost all configuration files.

Note: My dotfiles is [here](https://github.com/rhysd/dogfiles)


## Getting Started

1. Download [a released executable](https://github.com/rhysd/dotfiles/releases) and put it in `$PATH` or `$ go install github.com/rhysd/dotfiles`.
2. Change current directory to the directory you want to put a dotfiles repository.
3. Clone your dotfiles repository with `$ dotfiles clone`.
4. Enter the repository and run `$ dotfiles link --dry` to check which symlinks will be generated.
5. Write `.dotfiles/mappings.json` if needed.
6. `$ dotfiles link`
7. After you no longer need your configuration, remove all links with `$ dotfiles clean`.


## Usage

```
$ dotfiles {subcommand} [arguments]
```

### `clone` subcommand

Clone your dotfiles repository from remote.

```sh
# Clone git@github.com:rhysd/dotfiles.git into current directory
$ dotfiles clone rhysd

# Clone https://github.com/rhysd/dotfiles.git into current directory
$ dotfiles clone rhysd --https

# You can explicitly specify the repository name
$ dotfiles clone rhysd/dogfiles

# You can also use full-path
$ dotfiles clone git@bitbucket.org:rhysd/dotfiles.git
$ dotfiles clone https://your.site.com/dotfiles.git
```

### `link` subcommand

Set symbolic links to put your configuration files into proper places.

```sh
$ dotfiles link [options] [files...]
```

You can dry-run this command with `--dry` option.

If some `files` in dotfiles repository are specified, only they will be linked.

### `list` subcommand

Show all links set by this command.

```sh
$ dotfiles list
```

### `clean` subcommand

Remove all symbolic link put by `dotfiles link`.

```sh
$ dotfiles clean
```

### `update` subcommand

`git pull` your dotfiles repository from anywhere.

```sh
$ dotfiles update
```

### `selfupdate` subcommand

Update `dotfiles` binary (or `dotfiles.exe` on Windows) itself.

```sh
$ dotfiles selfupdate
```

## Default Mappings

It depends on your platform. Please see [source code](src/mappings.go).

## Symbolic Link Mappings

`dotfiles` command has sensible default mappings from configuration files in dotfiles repository to symbolic links put by `dotfiles link`.  And you can flexibly specify the mappings for your dotfiles manner.  Please create a `.dotfiles` directory and put a `.dotfiles/mappings.json` file in the root of your dotfiles repository.

Below is an example of `mappings.json`.  You can use `~` to represent a home directory.  As key, you can specify a name of file or directory in your dotfiles repository.  They will be linked to the corresponding values as symbolic links.

```json
{
  "gitignore": "~/.global.gitignore",
  "cabal_config": "~/.cabal/config"
}
```

In addition, you can define platform specific mappings with below mappings JSON files.

- `.dotfiles/mappings_unixlike.json`: Will link the mappings in Linux or macOS.
- `.dotfiles/mappings_linux.json`: Will link the mappings in Linux.
- `.dotfiles/mappings_darwin.json`: Will link the mappings in macOS.
- `.dotfiles/mappings_windows.json`: Will link the mappings in Windows.

Below is an example of `.dotfiles/mappings_darwin.json`.

```json
{
  "keyremap4macbook.xml": "~/Library/Application Support/Karabiner/private.xml",
  "mac.vimrc": "~/.mac.vimrc"
}
```

Values of the mappings object are basically strings representing destination paths, but they also can be arrays of strings. In the case, multiple symbolic links will be created for the source file.

For example, the following configuration will make two symbolic links `~/.vimrc` and `~/.config/nvim/init.vim` for `vimrc` source file.


```json
{
  "vimrc": ["~/.vimrc", "~/.config/nvim/init.vim"]
}
```

Real world example is [my dotfiles](https://github.com/rhysd/dogfiles/tree/master/.dotfiles).

## License

Licensed under [the MIT license](LICENSE.txt).

