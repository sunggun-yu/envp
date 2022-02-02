package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Proxy-Wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("prw %s, commit %s, built at %s", version, commit, date)
	},
}
