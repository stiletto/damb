package repo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

type DambDockerfileOptions struct {
	Context string
}

var reContext = regexp.MustCompile(`^#\s*DAMB:ctx:\s*(.*)\s*$`)

func ParseDambOptions(r io.Reader) (DambDockerfileOptions, error) {
	rdr := bufio.NewReader(r)
	result := DambDockerfileOptions{}
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return result, err
			}
		}
		if groups := reContext.FindStringSubmatch(line); groups != nil {
			if result.Context != "" {
				result.Context += " "
			}
			result.Context += groups[1]
		}
	}
	return result, nil
}

func LoadDambOptions(fname string) (DambDockerfileOptions, error) {
	f, err := os.Open(fname)
	result := DambDockerfileOptions{}
	if err != nil {
		return result, err
	}
	defer f.Close()
	result, err = ParseDambOptions(f)
	if err != nil {
		return result, fmt.Errorf("%s: %w", fname, err)
	}
	return result, nil
}
