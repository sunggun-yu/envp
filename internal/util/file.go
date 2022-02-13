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
