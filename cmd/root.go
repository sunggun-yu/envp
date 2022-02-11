package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
	"github.com/sunggun-yu/envp/internal/shell"
)

const (
	ConfigKeyDefaultProfile = "default"
	ConfigKeyProfile        = "profiles" // viper sub section key for profile
)

var (
	Config  config.Config
	rootCmd = rootCommand()
)

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

	cmd := &cobra.Command{
		Use:               "envp profile-name [flags] -- [command line to execute, e.g. kubectl]",
		Short:             "ENVP is cli wrapper that sets environment variables by profile when you execute the command line",
		SilenceUsage:      true,
		Example:           cmdExampleRoot(),
		ValidArgsFunction: ValidArgsProfileList,
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

			var profile *config.Profile
			var command []string
			var err error

			/*
				envp -- command         : ArgsLenAtDash == 0
				envp profile -- command : ArgsLenAtDash == 1
			*/
			switch {
			case cmd.ArgsLenAtDash() == 0:
				profile, err = Config.DefaultProfile()
				command = args
			case cmd.ArgsLenAtDash() == 1:
				profile, err = Config.Profiles.FindProfile(args[0])
				command = args[1:]
			}
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
			}

			// Execute command
			if err := shell.Execute(command, profile.Env); err != nil {
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
	viper.SetDefault("profiles", config.Profiles{})
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

	// Init Profiles
	// should watch write config event
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// reload profiles, and profile key list
		initProfiles()
	})

	// write config file with current config that is readed
	// this write will be helpful for the case config file is existing but empty
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

// unmarshal config.Config
func initProfiles() {
	err := viper.Unmarshal(&Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
