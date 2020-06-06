package dambfile

import (
	"sync"

	"github.com/stiletto/damb/utils"
)

type Arg struct {
	Cmd []string `yaml:"cmd"`
	Val string   `yaml:"val"`
	mu  sync.Mutex
}

type Dambfile struct {
	ImagePrefix string          `yaml:"root"`
	Dockerfile  string          `yaml:"dockerfile"`
	RefNameCmd  []string        `yaml:"refname_cmd"`
	BuildCmd    []string        `yaml:"build_cmd"`
	Args        map[string]*Arg `yaml:"args"`

	DirRoot string `yaml:"-"`
	RefName string `yaml:"-"`
}

func DefaultDambfile() *Dambfile {
	return &Dambfile{
		Dockerfile:  "Dockerfile",
		ImagePrefix: "local/",
		RefNameCmd:  []string{"git", "symbolic-ref", "--short", "HEAD"},
		BuildCmd:    []string{"docker", "build"},
	}
}

func (a *Arg) Get() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.Cmd) == 0 {
		return a.Val, nil
	}
	val, err := utils.CommandOutput(a.Cmd...)
	if err != nil {
		return "", err
	}
	a.Val = val
	a.Cmd = nil
	return val, nil
}
