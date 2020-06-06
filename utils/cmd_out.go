package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CommandOutput(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Command is empty")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = nil
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}
