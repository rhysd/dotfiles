package dotfiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

type Mappings map[string]AbsolutePath
type MappingsJson map[string]string

var DefaultMappings = map[string]MappingsJson{
	"windows": MappingsJson{
		".gvimrc": "~/vimfiles/gvimrc",
		".vim":    "~/vimfiles",
		".vimrc":  "~/vimfiles/vimrc",
	},
	"linux": MappingsJson{
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
	"darwin": MappingsJson{
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

func parseMappingsJson(file string) (MappingsJson, error) {
	var m MappingsJson

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

func convertMappingsJsonToMappings(json MappingsJson) (Mappings, error) {
	if json == nil {
		return nil, nil
	}
	m := make(Mappings, len(json))
	for k, v := range json {
		if k == "" {
			return nil, fmt.Errorf("Empty key cannot be included.  Note: Corresponding value is '%s'", v)
		}
		p, err := NewAbsolutePath(v)
		if err != nil {
			return nil, err
		}
		m[k] = p
	}
	return m, nil
}

func mergeMappingsFromFile(dist *Mappings, file string) error {
	j, err := parseMappingsJson(file)
	if err != nil {
		return err
	}
	if j == nil {
		return nil
	}

	m, err := convertMappingsJsonToMappings(j)
	if err != nil {
		return err
	}

	for k, v := range m {
		(*dist)[k] = v
	}

	return nil
}

func GetMappingsForPlatform(platform, parent string) (Mappings, error) {
	m, err := convertMappingsJsonToMappings(DefaultMappings[platform])
	if err != nil {
		return nil, err
	}
	if m == nil {
		m = Mappings{}
	}

	if err := mergeMappingsFromFile(&m, path.Join(parent, "mappings.json")); err != nil {
		return nil, err
	}

	if err := mergeMappingsFromFile(&m, path.Join(parent, fmt.Sprintf("mappings_%s.json", platform))); err != nil {
		return nil, err
	}

	// TODO:
	// Normalize path (e.g. ~/.foo -> /path/to/home/.foo)

	// TODO:
	// Validate mappings: linked path must be absolute

	return m, nil
}

func GetMappings(config_dir string) (Mappings, error) {
	return GetMappingsForPlatform(runtime.GOOS, config_dir)
}
