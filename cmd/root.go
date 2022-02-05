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
	cliName                 = "envp"
	ConfigKeyDefaultProfile = "default"
	ConfigKeyProfile        = "profiles" // viper sub section key for profile
)

// profile name from flag or config section "use"
var profile string

// unmarshalled object from selected profile in the config file
var currentProfile config.Profile

// root command that perform the command execution
var rootCmd = &cobra.Command{
	Use:          fmt.Sprintf("%v [flags] -- [command line to execute, such like kubectl]", cliName),
	Short:        fmt.Sprintf("%v is command line wrapper with selected env var profile", cliName),
	SilenceUsage: true,
	// TODO: externalize/refactoring
	Example: `
  # run command with selected profile. assuming HTTPS_PROXY is set in the profile
  envp use profile
  envp -- kubectl cluster-info
  envp -- kubectl get pods

  # specify env var profile to use
  envp --profile my-proxy -- kubectl get namespaces
  envp -p my-proxy -- kubectl get pods
	`,
	Args: cobra.MatchAll(
		cobra.MinimumNArgs(1),
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
		// default profile that was selected by "use" command will be used if profile flag is omitted
		if profile == "" && viper.GetString(ConfigKeyDefaultProfile) != "" {
			profile = viper.GetString(ConfigKeyDefaultProfile)
		}
		// validate if selected profile is existing in the config
		selected := viper.Sub(ConfigKeyProfile).Sub(profile)

		// check if selected profile from flag or use section in the config is existing
		if selected == nil {
			return fmt.Errorf("profile '%v' is not existing", profile)
		}

		// unmarshall to Profile
		err := selected.Unmarshal(&currentProfile)
		if err != nil {
			return fmt.Errorf("profile '%v' malformed configuration %e", profile, err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// first arg should be the command to execute
		command := args[0]

		// check if command can be found in the PATH
		binary, err := exec.LookPath(command)
		if err != nil {
			return err
		}

		// set additional environment variables to the session
		for _, e := range currentProfile.Env {
			os.Setenv(e.Name, e.Value)
		}

		// run commmand
		if err := syscall.Exec(binary, args, os.Environ()); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	// set optional flag "profile". so that user can select profile without swithing
	// selected profile by `use` command should be the profile if it is omitted
	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "usage string")
}

func initConfig() {

	// read config file
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// configPath := filepath.Join(home, "/.config/envp/config.toml")
	configPath := filepath.Join(home, ".config/envp")

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
