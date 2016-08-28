`dotfiles` Command
==================

**UNDER DEVELOPMENT. PLEASE DO NOT USE THIS YET.**

This repository provides `dotfiles` command to manage your [dotfiles](http://dotfiles.github.io/).  It manages your dotfiles repository and symbolic links to use the configurations.

This command has below goals:

- **One binary executable**: If you want to set configuration files in a remote server, all you have to do is sending a binary to the remote.
- **Do one thing and to it well**: This command manages only a dotfiles repository.  Does not handle any other dependencies.  If you want full-setup including dependencies, you should use more suitable tool such as [Ansible](https://www.ansible.com/).  And then use `dotfiles` command from it.
- **Less dependency**: Only depends on `git` command.
- **Sensible defaults**: Many sensible default symbolic link mappings are pre-defined.  You need not to specify the mappings for almost all configuration files.


## Installation

TODO


## Usage

```
$ dotfiles {subcommand} [arguments]
```

### `clone` subcommand

Clone your dotfiles repository from remote.

```sh
# Clone git@github.com:rhysd/dotfiles.git into current directory
$ dotfiles clone rhysd

# You can explicitly specify the repository name
$ dotfiles clone rhysd/dogfiles

# You can also use full-path
$ dotfiles clone git@bitbucket.org:rhysd/dotfiles.git
$ dotfiles clone https://your.site.com/dotfiles.git
```

### `link` subcommand

Set symbolic links to put your configuration files into proper places.

```sh
$ dotfiles link [options]
```

You can dry-run this command with `--dry` option.

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

## Default Mappings

TODO


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

- `.dotfiles/mappings_linux.json`: Will link the mappings in Linux.
- `.dotfiles/mappings_mac.json`: Will link the mappings in macOS.
- `.dotfiles/mappings_windows.json`: Will link the mappings in Windows.

Below is an example of `.dotfiles/mappings_mac.json`.

```json
{
  "keyremap4macbook.xml": "~/Library/Application Support/Karabiner/private.xml",
  "mac.vimrc": "~/.mac.vimrc"
}
```


## License

Licensed under [the MIT license](LICENSE.txt).

