package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
)

func FindUp(start string, names []string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		for _, name := range names {
			candidate := filepath.Join(dir, name)
			if err := syscall.Access(candidate, syscall.F_OK); err == nil {
				return candidate, nil
			}
		}
		nextDir := filepath.Dir(dir)
		if dir == nextDir {
			fnames := strings.Join(names, " or ")
			return "", fmt.Errorf("Couldn't find %s in current directory and its parent directories", fnames)
		}
		dir = nextDir
	}

}
