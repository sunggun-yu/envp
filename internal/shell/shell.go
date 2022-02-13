package shell

import (
	"os"
	"os/exec"

	"github.com/sunggun-yu/envp/internal/config"
	"github.com/sunggun-yu/envp/internal/util"
)

// TODO: refactoring, cleanup
// TODO: considering of using context
// TODO: poc of using forkExec and handling sigs, norifying sigs via channel and so on.

// Execute executes given command
func Execute(cmd []string, env config.Envs) error {
	return execCommand(cmd[0], cmd, env)
}

// StartShell runs default shell of user to create new shell session
func StartShell(env config.Envs) error {
	sh := os.Getenv("SHELL")
	return execCommand(sh, []string{sh, "-c", sh}, env)
}

// execCommand executes the os/exec Command with environment variales injection
func execCommand(argv0 string, argv []string, env config.Envs) error {
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

// setEnvs sets env vars of profile to system
// also, it expand abs path and set it as value of env var if the value start with "~" or "$HOME".
func setEnvs(env config.Envs) {
	// inject environment variables of profile
	for _, e := range env {
		// it's ok to ignore error. it returns original value if it doesn't contain the home path
		v, _ := util.ExpandHomeDir(e.Value)
		os.Setenv(e.Name, v)
	}
}
