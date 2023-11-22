package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/sunggun-yu/envp/internal/config"
)

// TODO: refactoring, cleanup
// TODO: considering of using context
// TODO: poc of using forkExec and handling sigs, notifying sigs via channel and so on.

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
func (s *ShellCommand) Execute(cmd []string, profile *config.NamedProfile) error {
	return s.execCommand(cmd[0], cmd, profile)
}

// StartShell runs default shell of user to create new shell session
func (s *ShellCommand) StartShell(profile *config.NamedProfile) error {
	sh := os.Getenv("SHELL")

	// use /bin/sh if SHELL is not set
	if sh == "" {
		sh = "/bin/sh"
	}

	// TODO: do some template
	// print start of session message
	s.Stdout.Write([]byte(fmt.Sprintln(color.GreenString("Starting ENVP session..."), color.RedString(profile.Name))))
	s.Stdout.Write([]byte(fmt.Sprintln("> press ctrl+d or type exit to close session")))

	// execute the command
	err := s.execCommand(sh, []string{sh, "-c", sh}, profile)
	if err != nil {
		s.Stderr.Write([]byte(fmt.Sprintln(color.MagentaString(err.Error()))))
	}

	// TODO: do some template
	// print end of session message
	s.Stdout.Write([]byte(fmt.Sprintln(color.GreenString("ENVP session closed..."), color.RedString(profile.Name))))

	return err
}

// execCommand executes the os/exec Command with environment variables injection
func (s *ShellCommand) execCommand(argv0 string, argv []string, profile *config.NamedProfile) error {
	// first arg should be the command to execute
	// check if command can be found in the PATH
	binary, err := exec.LookPath(argv0)
	if err != nil {
		return err
	}

	// TODO: refactor and make this clear separation of parsing and command substitution
	err = parseEnvs(profile.Env)
	if err != nil {
		return err
	}

	// create command for binary
	cmd := s.createCommand(&profile.Env, binary)
	// set args
	cmd.Args = argv
	// set ENVP_PROFILE
	cmd.Env = appendEnvpProfile(cmd.Env, profile.Name)

	// run init-script
	if err := s.executeInitScript(profile); err != nil {
		return err
	}

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// executeInitScript executes the initial script for the shell
func (s *ShellCommand) executeInitScript(profile *config.NamedProfile) error {
	// Return if profile or init-script is empty
	if profile == nil || profile.InitScript == nil {
		return nil
	}

	// loop and run init script in order
	for _, initScript := range profile.InitScripts() {
		cmd := s.createCommand(&profile.Env, "/bin/sh", "-c", initScript)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("init-script error: %w", err)
		}
	}
	return nil
}

// createCommand creates an *exec.Cmd instance configured with the provided command, arguments,
// environment variables, and associates Stdin, Stdout, and Stderr with the ShellCommand instance.
func (s *ShellCommand) createCommand(envs *config.Envs, cmd string, arg ...string) *exec.Cmd {

	c := exec.Command(cmd, arg...)
	c.Stdin = s.Stdin
	c.Stdout = s.Stdout
	c.Stderr = s.Stderr

	// init command Env with os.Environ()
	c.Env = os.Environ()
	// append config Envs to command
	c.Env = append(c.Env, envs.Strings()...)

	return c
}

// parseEnvs parse Env values with shell echo
func parseEnvs(envs config.Envs) (errs error) {
	for _, e := range envs {
		// parse env value with shell echo
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("echo %s", e.Value))
		// append envs to cmd that runs command substitution as well to support the case that reuse env var as ref with substitution
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envs.Strings()...)

		// it never occurs error since it is processed with shell echo.
		// so that, it will not exit 1 even command substitution has error. and just print out empty line when it errors
		output, _ := cmd.Output()
		// trim new lines from result
		result := strings.TrimRight(string(output), "\r\n")
		// use os.ExpandEnv to replace all the ${var} or $var in the string according to the values of the current environment variables.
		// so that $HOME will be replaced to current user's abs home dir
		result = os.ExpandEnv(result)

		if len(e.Value) > 0 && len(result) == 0 {
			// join errors
			errs = errors.Join(errs, fmt.Errorf("[envp] error processing value of %s: %s", e.Name, e.Value))
		} else {
			e.Value = result
		}
	}
	return errs
}

// appendEnvpProfile set ENVP_PROFILE env var to leverage profile info in the shell prompt, such as starship.
func appendEnvpProfile(envs []string, profile string) []string {
	envs = append(envs, fmt.Sprintf("%s=%s", envpEnvVarKey, profile))
	return envs
}
