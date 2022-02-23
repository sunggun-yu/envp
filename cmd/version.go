package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCommand())
}

func versionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of envp",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(rootCmd.Version)
		},
	}
	return cmd
}
