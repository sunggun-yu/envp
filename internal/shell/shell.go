package shell

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sunggun-yu/envp/internal/config"
)

// TODO: refactoring, cleanup

// ExecCmd execute command
func Execute(cmd []string, env []config.Env) error {
	return ExecCommand(cmd[0], cmd, env)
}

// StartShell runs default shell of user to create new shell session
func StartShell(env []config.Env) error {
	sh := os.Getenv("SHELL")

	if err := ExecCommand(sh, []string{sh}, env); err != nil {
		return err
	}
	return nil
}

// TODO: make it private once evaluation between exec.Command and syscall.Exec done
// ExecCommand executes the os/exec Command with environment variales injection
func ExecCommand(argv0 string, argv []string, env []config.Env) error {
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
	cmd.Env = os.Environ()

	// run commmand
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// TODO: make it private once evaluation between exec.Command and syscall.Exec done
// execute the command with syscall
func ExecuteWithSyscall(argv0 string, argv []string, env []config.Env) error {
	// first arg should be the command to execute
	// check if command can be found in the PATH
	binary, err := exec.LookPath(argv0)
	if err != nil {
		return err
	}

	// set environment variables
	setEnvs(env)

	// run commmand
	if err := syscall.Exec(binary, argv, os.Environ()); err != nil {
		return err
	}

	return nil
}

// set env vars of profile to system
func setEnvs(env []config.Env) {
	// inject environment variables of profile
	for _, e := range env {
		os.Setenv(e.Name, e.Value)
	}
}
