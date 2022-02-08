package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

func init() {
	rootCmd.AddCommand(listCommand())
}

// example of edit command
func cmdExampleList() string {
	return `
  envp list
  envp ls
  `
}

// listCommand prints out list of environment variable profiles
func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List all environment variable profiles",
		Aliases:      []string{"ls"},
		SilenceUsage: true,
		Example:      cmdExampleList(),
		RunE: func(cmd *cobra.Command, args []string) error {

			// current default profile name to compare
			defaultProfile := viper.GetString(ConfigKeyDefaultProfile)
			// print profiles. mark default profile with *
			for _, p := range profileList {
				if p == defaultProfile {
					fmt.Println("*", p)
				} else {
					fmt.Println(" ", p)
				}
			}
			return nil
		},
	}
	return cmd
}

// list all the profiles in dot "." format. e.g. mygroup.my-subgroup.my-profile
// Do DFS to build viper keys for profiles
func listProfileKeys(key string, profiles config.Profiles, arr *[]string) *[]string {
	for k, v := range profiles {
		var s string
		if key == "" {
			s = k
		} else {
			s = fmt.Sprint(key, ".", k)
		}
		// only Profile item has env items will be considered as profile
		// even group(parent Profile that has children Profiles) will be considered as Profile if it has env items.
		if len(v.Env) > 0 {
			*arr = append(*arr, s)
		}
		// recursion
		listProfileKeys(s, v.Profiles, arr)
	}
	return arr
}
