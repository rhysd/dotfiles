package dotfiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	return m, nil
}

func GetMappings(config_dir string) (Mappings, error) {
	return GetMappingsForPlatform(runtime.GOOS, config_dir)
}

func fileExists(file string) bool {
	s, err := os.Stat(file)
	return err == nil && !s.IsDir()
}

func link(from string, to AbsolutePath, dry bool) error {
	if to.IsEmpty() {
		// Note: Ignore if dist is specified 'null' in JSON
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	p := path.Join(cwd, from)
	if _, err := os.Stat(p); err != nil {
		if from[0] != '.' {
			return fmt.Errorf("'%s' does not exist. Please check the file in your dotfiles.", from)
		}

		p = path.Join(cwd, from[1:]) // Note: Omit '.'
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("Both '%s' and '%s' don't exist.  Please check the file in your dotfiles", from, from[1:])
		}
	}

	if _, err := os.Stat(string(to)); err == nil {
		fmt.Printf("'%s' already exists.  Skipped.\n", to)
		return nil
	}

	if err := os.MkdirAll(path.Dir(string(to)), os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	fmt.Printf("Link: '%s' -> '%s'\n", from, to)

	if dry {
		return nil
	}

	if err := os.Symlink(p, string(to)); err != nil {
		return err
	}

	return nil
}

func (mappings Mappings) CreateAllLinks(dry bool) error {
	for from, to := range mappings {
		if err := link(from, to, dry); err != nil {
			return err
		}
	}
	return nil
}

func (mappings Mappings) CreateSomeLinks(specified []string, dry bool) error {
	for _, from := range specified {
		if to, ok := mappings[from]; ok {
			if err := link(from, to, dry); err != nil {
				return err
			}
		}
	}
	return nil
}
