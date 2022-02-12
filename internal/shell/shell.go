package shell

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sunggun-yu/envp/internal/config"
)

// TODO: refactoring, cleanup
// TODO: considering of using context
// TODO: poc of using forkExec and handling sigs, norifying sigs via channel and so on.

// Execute executes given command
func Execute(cmd []string, env config.Envs) error {
	return ExecCommand(cmd[0], cmd, env)
}

// StartShell runs default shell of user to create new shell session
func StartShell(env config.Envs) error {
	sh := os.Getenv("SHELL")

	if err := ExecCommand(sh, []string{sh, "-c", sh}, env); err != nil {
		return err
	}
	return nil
}

// ExecCommand executes the os/exec Command with environment variales injection
// TODO: make it private once evaluation between exec.Command and syscall.Exec done
func ExecCommand(argv0 string, argv []string, env config.Envs) error {
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
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	// inject environment variables of profile
	setEnvs(env)
	// set environment variables to command
	// TODO: remove: passing env to cmd is not necessary in actually since setEnvs sets env vars to process
	cmd.Env = os.Environ()

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// ExecuteWithSyscall executes the command with syscall
// TODO: make it private once evaluation between exec.Command and syscall.Exec done
func ExecuteWithSyscall(argv0 string, argv []string, env config.Envs) error {
	// first arg should be the command to execute
	// check if command can be found in the PATH
	binary, err := exec.LookPath(argv0)
	if err != nil {
		return err
	}

	// set environment variables
	setEnvs(env)

	// run command
	if err := syscall.Exec(binary, argv, os.Environ()); err != nil {
		return err
	}

	return nil
}

// set env vars of profile to system
func setEnvs(env config.Envs) {
	// inject environment variables of profile
	for _, e := range env {
		os.Setenv(e.Name, e.Value)
	}
}
