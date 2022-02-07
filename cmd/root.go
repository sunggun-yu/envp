package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

const (
	ConfigKeyDefaultProfile = "default"
	ConfigKeyProfile        = "profiles" // viper sub section key for profile
)

var rootCmd = rootCommand()

// Execute execute the root command and sub commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
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
func rootCommand() *cobra.Command {
	// profile name from flag or config section "use"
	var profile string
	var command []string

	// unmarshalled object from selected profile in the config file
	var currentProfile config.Profile

	cmd := &cobra.Command{
		Use:          "envp profile-name [flags] -- [command line to execute, e.g. kubectl]",
		Short:        "ENVP is cli wrapper that sets environment variables by profile based configuration when you execute the command line",
		SilenceUsage: true,
		Example:      cmdExampleRoot(),
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
		// profile validation will be performed
		// global var for profile will be unmarshalled
		PreRunE: func(cmd *cobra.Command, args []string) error {
			/*
				envp -- command         : ArgsLenAtDash == 0
				envp profile -- command : ArgsLenAtDash == 1
			*/
			switch {
			case cmd.ArgsLenAtDash() == 0:
				// this case requires default profile.
				if viper.GetString(ConfigKeyDefaultProfile) == "" {
					printExample(cmd)
					return fmt.Errorf("default profile is not set. please set default profile")
				}

				profile = viper.GetString(ConfigKeyDefaultProfile)
				command = args
			case cmd.ArgsLenAtDash() == 1:
				profile = args[0]
				command = args[1:]
			}

			// check if selected profile is existing
			if viper.Sub(ConfigKeyProfile).Sub(profile) == nil {
				return fmt.Errorf("profile '%v' is not existing", profile)
			}

			// validate if selected profile is existing in the config
			selected := viper.Sub(ConfigKeyProfile).Sub(profile)
			// unmarshal to Profile
			err := selected.Unmarshal(&currentProfile)
			if err != nil {
				return fmt.Errorf("profile '%v' malformed configuration %e", profile, err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// first arg should be the command to execute
			// check if command can be found in the PATH
			binary, err := exec.LookPath(command[0])
			if err != nil {
				return err
			}

			// set environment variables to the session
			for _, e := range currentProfile.Env {
				os.Setenv(e.Name, e.Value)
			}

			// run commmand
			if err := syscall.Exec(binary, command, os.Environ()); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

// initConfig read the config file and initialize if config file is not existing
func initConfig() {
	// set default empty profile name
	viper.SetDefault("default", "")
	// set default empty profiles
	viper.SetDefault("profiles", config.Profile{})
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath(".config/envp")) // $HOME/.config/envp
	// write config file if file does not existing
	viper.SafeWriteConfig()

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	// write config file with current config that is readed
	// this write will be helpful for the case config file is existing but empty
	viper.WriteConfig()
}

// get config path. mkdir -p it not exist
func configPath(base string) string {
	// get $HOME
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// get config path : $HOME/.config/envp
	path := filepath.Join(home, base)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// mkdir -p if directory is not existing
			os.MkdirAll(path, 0755)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return path
}

// print example only
func printExample(cmd *cobra.Command) {
	fmt.Println("Example:")
	fmt.Println(cmd.Example)
}
