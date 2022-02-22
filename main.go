package main

import (
	"fmt"

	"github.com/sunggun-yu/envp/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// set version
	cmd.SetVersion(Version())
	cmd.Execute()
}

// Version returns version and build information. it will be injected from ldflags(goreleaser)
func Version() string {
	return fmt.Sprintf("envp %s, commit %s, built at %s", version, commit, date)
}
