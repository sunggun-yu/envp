package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/sunggun-yu/envp/internal/config"
	"github.com/sunggun-yu/envp/internal/util"
)

// TODO: refactoring, cleanup
// TODO: considering of using context
// TODO: poc of using forkExec and handling sigs, norifying sigs via channel and so on.

const envpEnvVarKey = "ENVP_PROFILE"

// ShellCommand is struct of shell command.
type ShellCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewShellCommand create ShellCommand with os stdin, stdout, and stderr as default
func NewShellCommand() *ShellCommand {
	return &ShellCommand{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Execute executes given command
func (s *ShellCommand) Execute(cmd []string, env config.Envs, profile string) error {
	return s.execCommand(cmd[0], cmd, env, profile)
}

// StartShell runs default shell of user to create new shell session
func (s *ShellCommand) StartShell(env config.Envs, profile string) error {
	sh := os.Getenv("SHELL")

	// TODO: do some template
	// print start of session message
	s.Stdout.Write([]byte(fmt.Sprintln(color.GreenString("Starting ENVP session..."), color.RedString(profile))))
	s.Stdout.Write([]byte(fmt.Sprintln(color.CyanString(env.String()))))
	s.Stdout.Write([]byte(fmt.Sprintln("> press ctrl+d or type exit to close session")))

	// execute the command
	err := s.execCommand(sh, []string{sh, "-c", sh}, env, profile)

	// TODO: do some template
	// print end of session message
	s.Stdout.Write([]byte(fmt.Sprintln(color.GreenString("ENVP session closed..."), color.RedString(profile))))

	return err
}

// execCommand executes the os/exec Command with environment variales injection
func (s *ShellCommand) execCommand(argv0 string, argv []string, env config.Envs, profile string) error {
	// first arg should be the command to execute
	// check if command can be found in the PATH
	binary, err := exec.LookPath(argv0)
	if err != nil {
		return err
	}

	// create command for binary
	cmd := exec.Command(binary)
	// set args
	cmd.Args = argv
	cmd.Stdout = s.Stdout
	cmd.Stdin = s.Stdin
	cmd.Stderr = s.Stderr
	// merge into os environment variables and set into the cmd
	cmd.Env = append(os.Environ(), appendEnvpProfile(parseEnvs(env), profile)...)

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// parseEnvs parse config.Envs to "VAR=VAL" format string slice
func parseEnvs(env config.Envs) []string {
	ev := []string{}
	for _, e := range env {
		// it's ok to ignore error. it returns original value if it doesn't contain the home path
		v, _ := util.ExpandHomeDir(e.Value)
		// Env.String() would not work for this case since we want to cover expanding the home dir path
		ev = append(ev, fmt.Sprintf("%s=%s", e.Name, v))
	}
	return ev
}

// appendEnvpProfile set ENVP_PROFILE env var to leverage profile info in the shell prompt, such as starship.
func appendEnvpProfile(envs []string, profile string) []string {
	envs = append(envs, fmt.Sprintf("%s=%s", envpEnvVarKey, profile))
	return envs
}
