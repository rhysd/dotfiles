package dotfiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/rhysd/abspath"
)

type NothingLinkedError struct {
	RepoPath string
}

func (err NothingLinkedError) Error() string {
	if err.RepoPath == "" {
		return "Nothing was linked."
	}
	return fmt.Sprintf("Nothing was linked. '%s' was specified as dotfiles repository. Please check it", err.RepoPath)
}

// UnixLikePlatformName is a special platform name used commonly for Unix-like platform
const UnixLikePlatformName = "unixlike"

type Mappings map[string][]abspath.AbsPath
type MappingsJSON map[string]interface{}

var DefaultMappings = map[string]MappingsJSON{
	"windows": MappingsJSON{
		".gvimrc": "~/vimfiles/gvimrc",
		".vim":    "~/vimfiles",
		".vimrc":  "~/vimfiles/vimrc",
	},
	UnixLikePlatformName: MappingsJSON{
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
		"agignore":       "~/.agignore",
		"bash_login":     "~/.bash_login",
		"bash_profile":   "~/.bash_profile",
		"bashrc":         "~/.bashrc",
		"emacs.d":        "~/.emacs.d",
		"emacs.el":       "~/.emacs.d/init.el",
		"eslintrc":       "~/.eslintrc",
		"eslintrc.json":  "~/.eslintrc.json",
		"eslintrc.yml":   "~/.eslintrc.yml",
		"gvimrc":         "~/.gvimrc",
		"npmrc":          "~/.npmrc",
		"profile":        "~/.profile",
		"pryrc":          "~/.pryrc",
		"pylintrc":       "~/.pylintrc",
		"tmux.conf":      "~/.tmux.conf",
		"vim":            "~/.vim",
		"vimrc":          "~/.vimrc",
		"zlogin":         "~/.zlogin",
		"zprofile":       "~/.zprofile",
		"zshenv":         "~/.zshenv",
		"zshrc":          "~/.zshrc",
		"init.el":        "~/.emacs.d/init.el",
		"peco":           "~/.config/peco",
	},
	"linux": MappingsJSON{
		".Xmodmap":    "~/.Xmodmap",
		".Xresources": "~/.Xresources",
		"Xmodmap":     "~/.Xmodmap",
		"Xresources":  "~/.Xresources",
		"rc.lua":      "~/.config/rc.lua",
	},
	"darwin": MappingsJSON{
		".htoprc": "~/.htoprc",
		"htoprc":  "~/.htoprc",
	},
}

func parseMappingsJSON(file abspath.AbsPath) (MappingsJSON, error) {
	var m MappingsJSON

	bytes, err := ioutil.ReadFile(file.String())
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

func expandPath(s string) (abspath.AbsPath, error) {
	if s[0] != '~' && s[0] != '/' {
		return abspath.AbsPath{}, fmt.Errorf("Value of mappings must be an absolute path like '/foo/.bar' or '~/.foo': %s", s)
	}
	return abspath.ExpandFromSlash(s)
}

func convertMappingsJSONToMappings(json MappingsJSON) (Mappings, error) {
	if json == nil {
		return nil, nil
	}
	m := make(Mappings, len(json))
	for k, v := range json {
		if k == "" {
			return nil, fmt.Errorf("Empty key cannot be included.  Note: Corresponding value is '%s'", v)
		}
		switch v := v.(type) {
		case string:
			if v == "" {
				// Note: Ignore if dist is specified 'null' in JSON
				continue
			}
			p, err := expandPath(v)
			if err != nil {
				return nil, err
			}
			m[k] = []abspath.AbsPath{p}
		case []interface{}:
			m[k] = make([]abspath.AbsPath, 0, len(v))
			for _, iface := range v {
				s, ok := iface.(string)
				if !ok {
					return nil, fmt.Errorf("Type of value must be string or string[]: %v", v)
				}
				p, err := expandPath(s)
				if err != nil {
					return nil, err
				}
				m[k] = append(m[k], p)
			}
		default:
			return nil, fmt.Errorf("Value of mappings object must be string or string[]: %v", v)
		}
	}
	return m, nil
}

func mergeMappingsFromDefault(dist Mappings, platform string) error {
	m, err := convertMappingsJSONToMappings(DefaultMappings[platform])
	if err != nil {
		return err
	}

	for k, v := range m {
		dist[k] = v
	}

	return nil
}

func mergeMappingsFromFile(dist Mappings, file abspath.AbsPath) error {
	j, err := parseMappingsJSON(file)
	if err != nil {
		return err
	}
	if j == nil {
		return nil
	}

	m, err := convertMappingsJSONToMappings(j)
	if err != nil {
		return err
	}

	for k, v := range m {
		dist[k] = v
	}

	return nil
}

func isUnixLikePlatform(platform string) bool {
	return platform == "linux" || platform == "darwin"
}

func GetMappingsForPlatform(platform string, parent abspath.AbsPath) (Mappings, error) {
	m := Mappings{}

	if isUnixLikePlatform(platform) {
		if err := mergeMappingsFromDefault(m, UnixLikePlatformName); err != nil {
			return nil, err
		}
	}
	if err := mergeMappingsFromDefault(m, platform); err != nil {
		return nil, err
	}

	if err := mergeMappingsFromFile(m, parent.Join("mappings.json")); err != nil {
		return nil, err
	}

	if isUnixLikePlatform(platform) {
		if err := mergeMappingsFromFile(m, parent.Join(fmt.Sprintf("mappings_%s.json", UnixLikePlatformName))); err != nil {
			return nil, err
		}
	}
	if err := mergeMappingsFromFile(m, parent.Join(fmt.Sprintf("mappings_%s.json", platform))); err != nil {
		return nil, err
	}

	return m, nil
}

func GetMappings(configDir abspath.AbsPath) (Mappings, error) {
	return GetMappingsForPlatform(runtime.GOOS, configDir)
}

func fileExists(file string) bool {
	s, err := os.Stat(file)
	return err == nil && !s.IsDir()
}

func link(from string, to abspath.AbsPath, dry bool) (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	p := filepath.Join(cwd, from)
	if _, err := os.Stat(p); err != nil {
		return false, nil
	}

	if _, err := os.Stat(to.String()); err == nil {
		// Target already exists. Skipped.
		fmt.Printf("Exist: '%s' -> '%s'\n", from, to.String())
		return true, nil
	}

	if err := os.MkdirAll(to.Dir().String(), os.ModeDir|os.ModePerm); err != nil {
		return false, err
	}

	color.Cyan("Link:  '%s' -> '%s'\n", from, to.String())

	if dry {
		return true, nil
	}

	if err := os.Symlink(p, to.String()); err != nil {
		return false, err
	}

	return true, nil
}

func (mappings Mappings) CreateAllLinks(dry bool) error {
	count := 0
	for from, tos := range mappings {
		for _, to := range tos {
			linked, err := link(from, to, dry)
			if err != nil {
				return err
			}
			if linked {
				count++
			}
		}
	}

	if count == 0 {
		return &NothingLinkedError{}
	}

	return nil
}

func (mappings Mappings) CreateSomeLinks(specified []string, dry bool) error {
	count := 0
	for _, from := range specified {
		if tos, ok := mappings[from]; ok {
			for _, to := range tos {
				linked, err := link(from, to, dry)
				if err != nil {
					return err
				}
				if linked {
					count++
				}
			}
		}
	}

	if count == 0 && specified != nil && len(specified) > 0 {
		return &NothingLinkedError{}
	}

	return nil
}

func getLinkSource(repo, to abspath.AbsPath) (string, error) {
	s, err := os.Lstat(to.String())
	if err != nil {
		// Note: Symlink not found
		return "", nil
	}

	if s.Mode()&os.ModeSymlink != os.ModeSymlink {
		return "", nil
	}

	source, err := os.Readlink(to.String())
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(source, repo.String()) {
		// Note: When the symlink is not linked from dotfiles repository.
		return "", nil
	}

	return source, nil
}

func (mappings Mappings) unlink(repo, to abspath.AbsPath) (bool, error) {
	source, err := getLinkSource(repo, to)
	if source == "" || err != nil {
		return false, err
	}

	if err := os.Remove(to.String()); err != nil {
		return false, err
	}

	fmt.Printf("Unlink: '%s' -> '%s'\n", source, to.String())

	return true, nil
}

func (mappings Mappings) UnlinkAll(repo abspath.AbsPath) error {
	count := 0
	for _, tos := range mappings {
		for _, to := range tos {
			unlinked, err := mappings.unlink(repo, to)
			if err != nil {
				return err
			}
			if unlinked {
				count++
			}
		}
	}

	if count == 0 {
		fmt.Printf("No symlink was removed (dotfiles: '%s').\n", repo.String())
	}

	return nil
}

func (mappings Mappings) ActualLinks(repo abspath.AbsPath) (map[string]string, error) {
	ret := map[string]string{}
	for _, tos := range mappings {
		for _, to := range tos {
			s, err := getLinkSource(repo, to)
			if err != nil {
				return nil, err
			}
			if s != "" {
				ret[s] = to.String()
			}
		}
	}
	return ret, nil
}
