package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
	"github.com/sunggun-yu/envp/internal/shell"
)

var (
	configFile     *config.ConfigFile                     // ConfigFile instance that is shared across the sub-commands
	configFileName = "$HOME/.config/envp/config.yaml"     // config file path
	rootCmd        = rootCommand(shell.NewShellCommand()) // root command with default setup of shell command
)

// init
func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig initialize the config file
func initConfig() {
	if cfg, err := config.NewConfigFile(configFileName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		configFile = cfg
	}
}

// Execute execute the root command and sub commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// SetVersion set the version of command
func SetVersion(version string) {
	rootCmd.Version = version
}

// example of add command
func cmdExampleRoot() string {
	return `
  # run command with selected environment variable profile.
  # (example is assuming HTTPS_PROXY is set in the profile)
  envp use profile
  envp -- kubectl cluster-info
  envp -- kubectl get pods
  
  # specify env var profile to use
  envp profile-name -- kubectl get namespaces
  `
}

// rootCommand sets environment variable and execute command line
func rootCommand(sh *shell.ShellCommand) *cobra.Command {

	cmd := &cobra.Command{
		Use:               "envp profile-name [flags] -- [command line to execute, e.g. kubectl]",
		Short:             "ENVP is cli wrapper that sets environment variables by profile when you execute the command line",
		SilenceUsage:      true,
		Example:           cmdExampleRoot(),
		ValidArgsFunction: validArgsProfileList,
		Args: cobra.MatchAll(
			func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					cmd.Help()
					os.Exit(0)
				}
				return nil
			},
			func(cmd *cobra.Command, args []string) error {
				if len(args) > 0 && cmd.ArgsLenAtDash() < 0 {
					return fmt.Errorf("command should start after --")
				}
				return nil
			},
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			var cfg *config.Config
			var profile *config.NamedProfile
			var command []string
			var err error

			cfg, err = configFile.Read()
			if err != nil {
				return err
			}

			/*
				envp -- command         : ArgsLenAtDash == 0
				envp profile -- command : ArgsLenAtDash == 1
			*/
			switch {
			case cmd.ArgsLenAtDash() == 0:
				profile, err = cfg.DefaultProfile()
				command = args
			case cmd.ArgsLenAtDash() > 0:
				profile, err = cfg.Profile(args[0])
				// only args after double dash "--"" should be considered as command
				command = args[cmd.ArgsLenAtDash():]
			}
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
			}

			// Execute command
			if err := sh.Execute(command, profile); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
