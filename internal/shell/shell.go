package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/sunggun-yu/envp/internal/config"
	"github.com/sunggun-yu/envp/internal/util"
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

	// create command for binary
	cmd := exec.Command(binary)
	// set args
	cmd.Args = argv
	cmd.Stdout = s.Stdout
	cmd.Stdin = s.Stdin
	cmd.Stderr = s.Stderr

	// init cmd.Env with os.Environ()
	cmd.Env = os.Environ()
	// set ENVP_PROFILE
	cmd.Env = appendEnvpProfile(cmd.Env, profile.Name)

	err = parseEnvs(profile.Env)
	if err != nil {
		return err
	}

	// merge into os environment variables and set into the cmd
	cmd.Env = append(cmd.Env, profile.Env.Strings()...)

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

	// just return if init-script is empty
	if profile == nil || len(profile.InitScript) == 0 {
		return nil
	}

	initCmd := exec.Command("/bin/sh", "-c", profile.InitScript)
	initCmd.Stdin = s.Stdin
	initCmd.Stdout = s.Stdout
	initCmd.Stderr = s.Stderr

	initCmd.Env = os.Environ()
	// append envs to init-script so that both init and main shell can maintain common env var
	initCmd.Env = append(initCmd.Env, profile.Env.Strings()...)
	err := initCmd.Run()
	if err != nil {
		return fmt.Errorf("init-script error: %w", err)
	}
	return nil
}

// parseEnvs parse config.Envs to "VAR=VAL" format string slice
func parseEnvs(envs config.Envs) (errs error) {
	for _, e := range envs {
		// it's ok to ignore error. it returns original value if it doesn't contain the home path
		e.Value, _ = util.ExpandHomeDir(e.Value)
		// parse command substitution value like $(some-command). treat error to let user to know there is error with it
		v, err := processCommandSubstitutionValue(e.Value, envs)
		if err != nil {
			// join errors
			errs = errors.Join(errs, fmt.Errorf("[envp] error processing value of %s: %s", e.Name, err))
		} else {
			e.Value = v
		}
	}
	return errs
}

// appendEnvpProfile set ENVP_PROFILE env var to leverage profile info in the shell prompt, such as starship.
func appendEnvpProfile(envs []string, profile string) []string {
	envs = append(envs, fmt.Sprintf("%s=%s", envpEnvVarKey, profile))
	return envs
}

// processCommandSubstitutionValue checks whether the env value is in the format of shell substitution $() and runs the shell to replace the env value with the result of its execution.
func processCommandSubstitutionValue(val string, envs config.Envs) (string, error) {
	// check if val is pattern of command substitution using regex
	// support only $() substitution. not support `` substitution
	re := regexp.MustCompile(`^\$\((.*?)\)`) // use MustCompile. no expect it's failing

	matches := re.FindStringSubmatch(val)
	if len(matches) < 2 {
		// no valid script found. just return original value
		return val, nil
	}

	script := strings.TrimSpace(matches[1])
	cmd := exec.Command("sh", "-c", script)
	// append envs to cmd that runs command substitution as well to support the case that reuse env var as ref with substitution
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs.Strings()...)

	// output, err := cmd.CombinedOutput()
	output, err := cmd.Output()
	if err != nil {
		return val, fmt.Errorf("error executing script: %v", err)
	}

	return strings.TrimRight(string(output), "\r\n"), nil
}
