package dotfiles

import (
	"os"
	"path"
)

func Link(specified []string, dry bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	m, err := GetMappings(path.Join(cwd, ".dotfiles"))
	if err != nil {
		return err
	}

	if specified == nil || len(specified) == 0 {
		return m.CreateAllLinks(dry)
	} else {
		return m.CreateSomeLinks(specified, dry)
	}
}
