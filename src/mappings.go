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
type MappingsJSON map[string][]string

var DefaultMappings = map[string]MappingsJSON{
	"windows": MappingsJSON{
		".gvimrc": []string{"~/vimfiles/gvimrc"},
		".vim":    []string{"~/vimfiles"},
		".vimrc":  []string{"~/vimfiles/vimrc"},
	},
	UnixLikePlatformName: MappingsJSON{
		".agignore":      []string{"~/.agignore"},
		".bash_login":    []string{"~/.bash_login"},
		".bash_profile":  []string{"~/.bash_profile"},
		".bashrc":        []string{"~/.bashrc"},
		".emacs.d":       []string{"~/.emacs.d"},
		".emacs.el":      []string{"~/.emacs.d/init.el"},
		".eslintrc":      []string{"~/.eslintrc"},
		".eslintrc.json": []string{"~/.eslintrc.json"},
		".eslintrc.yml":  []string{"~/.eslintrc.yml"},
		".gvimrc":        []string{"~/.gvimrc"},
		".npmrc":         []string{"~/.npmrc"},
		".profile":       []string{"~/.profile"},
		".pryrc":         []string{"~/.pryrc"},
		".pylintrc":      []string{"~/.pylintrc"},
		".tmux.conf":     []string{"~/.tmux.conf"},
		".vim":           []string{"~/.vim"},
		".vimrc":         []string{"~/.vimrc"},
		".zlogin":        []string{"~/.zlogin"},
		".zprofile":      []string{"~/.zprofile"},
		".zshenv":        []string{"~/.zshenv"},
		".zshrc":         []string{"~/.zshrc"},
		"agignore":       []string{"~/.agignore"},
		"bash_login":     []string{"~/.bash_login"},
		"bash_profile":   []string{"~/.bash_profile"},
		"bashrc":         []string{"~/.bashrc"},
		"emacs.d":        []string{"~/.emacs.d"},
		"emacs.el":       []string{"~/.emacs.d/init.el"},
		"eslintrc":       []string{"~/.eslintrc"},
		"eslintrc.json":  []string{"~/.eslintrc.json"},
		"eslintrc.yml":   []string{"~/.eslintrc.yml"},
		"gvimrc":         []string{"~/.gvimrc"},
		"npmrc":          []string{"~/.npmrc"},
		"profile":        []string{"~/.profile"},
		"pryrc":          []string{"~/.pryrc"},
		"pylintrc":       []string{"~/.pylintrc"},
		"tmux.conf":      []string{"~/.tmux.conf"},
		"vim":            []string{"~/.vim"},
		"vimrc":          []string{"~/.vimrc"},
		"zlogin":         []string{"~/.zlogin"},
		"zprofile":       []string{"~/.zprofile"},
		"zshenv":         []string{"~/.zshenv"},
		"zshrc":          []string{"~/.zshrc"},
		"init.el":        []string{"~/.emacs.d/init.el"},
		"peco":           []string{"~/.config/peco"},
	},
	"linux": MappingsJSON{
		".Xmodmap":    []string{"~/.Xmodmap"},
		".Xresources": []string{"~/.Xresources"},
		"Xmodmap":     []string{"~/.Xmodmap"},
		"Xresources":  []string{"~/.Xresources"},
		"rc.lua":      []string{"~/.config/rc.lua"},
	},
	"darwin": MappingsJSON{
		".htoprc": []string{"~/.htoprc"},
		"htoprc":  []string{"~/.htoprc"},
	},
}

type PathLink struct {
	src, dst string
}

func parseMappingsJSON(file abspath.AbsPath) (MappingsJSON, error) {
	var m map[string]interface{}

	bytes, err := ioutil.ReadFile(file.String())
	if err != nil {
		// Note:
		// It's not an error that the file is not found
		return nil, nil
	}

	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	maps := make(MappingsJSON, len(m))
	for k, v := range m {
		switch v := v.(type) {
		case string:
			maps[k] = []string{v}
		case []interface{}:
			vs := make([]string, 0, len(v))
			for _, iface := range v {
				s, ok := iface.(string)
				if !ok {
					return nil, fmt.Errorf("Value of mappings object must be string or string[]: %v", v)
				}
				vs = append(vs, s)
			}
			maps[k] = vs
		}
	}

	return maps, nil
}

func convertMappingsJSONToMappings(json MappingsJSON) (Mappings, error) {
	if json == nil {
		return nil, nil
	}
	m := make(Mappings, len(json))
	for k, vs := range json {
		if k == "" {
			return nil, fmt.Errorf("Empty key cannot be included.  Note: Corresponding value is '%s'", vs)
		}
		ps := make([]abspath.AbsPath, 0, len(vs))
		for _, v := range vs {
			if v == "" {
				continue
			}
			if v[0] != '~' && v[0] != '/' {
				return nil, fmt.Errorf("Value of mappings must be an absolute path like '/foo/.bar' or '~/.foo': %s", v)
			}
			p, err := abspath.ExpandFromSlash(v)
			if err != nil {
				return nil, err
			}
			ps = append(ps, p)
		}
		m[k] = ps
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

func (maps Mappings) CreateAllLinks(dry bool) error {
	created := false
	for from, tos := range maps {
		for _, to := range tos {
			linked, err := link(from, to, dry)
			if err != nil {
				return err
			}
			if linked {
				created = true
			}
		}
	}

	if !created {
		return &NothingLinkedError{}
	}

	return nil
}

func (maps Mappings) CreateSomeLinks(specified []string, dry bool) error {
	created := false
	for _, from := range specified {
		if tos, ok := maps[from]; ok {
			for _, to := range tos {
				linked, err := link(from, to, dry)
				if err != nil {
					return err
				}
				if linked {
					created = true
				}
			}
		}
	}

	if !created && len(specified) > 0 {
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

func (maps Mappings) unlink(repo, to abspath.AbsPath) (bool, error) {
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

func (maps Mappings) UnlinkAll(repo abspath.AbsPath) error {
	count := 0
	for _, tos := range maps {
		for _, to := range tos {
			unlinked, err := maps.unlink(repo, to)
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

func (maps Mappings) ActualLinks(repo abspath.AbsPath) ([]PathLink, error) {
	ret := make([]PathLink, 0, len(maps))
	for _, tos := range maps {
		for _, to := range tos {
			s, err := getLinkSource(repo, to)
			if err != nil {
				return nil, err
			}
			if s != "" {
				ret = append(ret, PathLink{s, to.String()})
			}
		}
	}
	return ret, nil
}
