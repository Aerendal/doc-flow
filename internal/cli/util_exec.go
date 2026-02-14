package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runShell(command string) error {
	c := fmt.Sprintf("set -euo pipefail\n%s", command)
	return execCommand("/bin/bash", "-lc", c)
}
