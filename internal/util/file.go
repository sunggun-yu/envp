package util

import (
	"os"
	"strings"
)

// ExpandHomeDir returns abs path of filename if it starts with one of `~` or `$HOME`
// e.g.
//   ~/.config/envp/config.yaml -> /<Users|home>/someone/.config/envp/config.yaml
//   $HOME/.config/envp/config.yaml -> /<Users|home>/someone/.config/envp/config.yaml
//   /tmp/somepath -> /tmp/somepath
func ExpandHomeDir(f string) (string, error) {
	homes := []string{"~", "$HOME"}
	for _, h := range homes {
		if strings.HasPrefix(f, h) {
			userHome, err := os.UserHomeDir()
			if err != nil {
				return f, err
			}
			return strings.Replace(f, h, userHome, 1), nil
		}
	}
	return f, nil
}

// EnsureConfigFilePath ensure the config file path.
// it mkdir -p the file's directory if not exist
// also returns abs path of file name if it starts from `~` or `$HOME`
// path must be "dir" not "file"
func EnsureConfigFilePath(path string) (string, error) {
	// expand home dir
	f, err := ExpandHomeDir(path)
	if err != nil {
		return path, err
	}
	// ensure if file is existing
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(f, 0755)
			if err != nil {
				return f, err
			}
		} else {
			// return other errors
			return f, err
		}
	}
	return f, nil
}
