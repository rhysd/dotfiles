package dotfiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

type Mappings map[string]string

var DefaultMappings = map[string]Mappings{
	"windows": Mappings{
		".gvimrc": "~/vimfiles/gvimrc",
		".vim":    "~/vimfiles",
		".vimrc":  "~/vimfiles/vimrc",
	},
	"linux": Mappings{
		".Xmodmap":       "~/.Xmodmap",
		".Xresources":    "~/.Xresources",
		".agignore":      "~/.agignore",
		".bash_login":    "~/.bash_login",
		".bash_profile":  "~/.bash_profile",
		".bashrc":        "~/.bashrc",
		".emacs.d":       "~/.emacs.d",
		".emacs.el":      "~/.emacs.d/init.el",
		".eslintrc":      "~/.eslintrc",
		".eslintrc.json": "~/.eslintrc.json",
		".eslintrc.yml":  "~/.eslintrc.yml",
		".gvimrc":        "~/.gvimrc",
		".npmrc":         "~/.npmrc",
		".profile":       "~/.profile",
		".pryrc":         "~/.pryrc",
		".pylintrc":      "~/.pylintrc",
		".tmux.conf":     "~/.tmux.conf",
		".vim":           "~/.vim",
		".vimrc":         "~/.vimrc",
		".zlogin":        "~/.zlogin",
		".zprofile":      "~/.zprofile",
		".zshenv":        "~/.zshenv",
		".zshrc":         "~/.zshrc",
		"init.el":        "~/.emacs.d/init.el",
		"peco":           "~/.config/peco",
		"rc.lua":         "~/.config/rc.lua",
	},
	"darwin": Mappings{
		".agignore":      "~/.agignore",
		".bash_login":    "~/.bash_login",
		".bash_profile":  "~/.bash_profile",
		".bashrc":        "~/.bashrc",
		".emacs.d":       "~/.emacs.d",
		".emacs.el":      "~/.emacs.d/init.el",
		".eslintrc":      "~/.eslintrc",
		".eslintrc.json": "~/.eslintrc.json",
		".eslintrc.yml":  "~/.eslintrc.yml",
		".gvimrc":        "~/.gvimrc",
		".htoprc":        "~/.htoprc",
		".npmrc":         "~/.npmrc",
		".profile":       "~/.profile",
		".pryrc":         "~/.pryrc",
		".pylintrc":      "~/.pylintrc",
		".tmux.conf":     "~/.tmux.conf",
		".vim":           "~/.vim",
		".vimrc":         "~/.vimrc",
		".zlogin":        "~/.zlogin",
		".zprofile":      "~/.zprofile",
		".zshenv":        "~/.zshenv",
		".zshrc":         "~/.zshrc",
		"init.el":        "~/.emacs.d/init.el",
		"peco":           "~/.config/peco",
	},
}

func parseMappingsJson(file string) (Mappings, error) {
	var m Mappings

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		// Note:
		// It's not an error that the file is not found
		return nil, nil
	}

	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func mergeMappingsFromFile(dist *Mappings, file string) error {
	m, err := parseMappingsJson(file)
	if err != nil {
		return err
	}
	if m == nil {
		return nil
	}

	for k, v := range m {
		(*dist)[k] = v
	}
	return nil
}

func GetMappingsForPlatform(platform, parent string) (Mappings, error) {
	m := DefaultMappings[platform]
	if m == nil {
		m = Mappings{}
	}

	if err := mergeMappingsFromFile(&m, path.Join(parent, "mappings.json")); err != nil {
		return nil, err
	}

	if err := mergeMappingsFromFile(&m, path.Join(parent, fmt.Sprintf("mappings_%s.json", platform))); err != nil {
		return nil, err
	}

	return m, nil
}

func GetMappings(config_dir string) (Mappings, error) {
	return GetMappingsForPlatform(runtime.GOOS, config_dir)
}
