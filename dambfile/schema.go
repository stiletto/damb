package dambfile

import (
	"fmt"
	"sync"

	"github.com/stiletto/damb/utils"
)

type Arg struct {
	Cmd []string `yaml:"cmd"`
	Val string   `yaml:"val"`
	mu  sync.Mutex
}

type Dambfile struct {
	Dockerfile    string              `yaml:"dockerfile"`
	BuildCmd      []string            `yaml:"build_cmd"`
	Args          map[string]*Arg     `yaml:"args"`
	Aliases       map[string][]string `yaml:"aliases"`
	CacheFromTags []string            `yaml:"cache_from_tags"`

	DirRoot     string `yaml:"-" json:"-"`
	ImageTag    string `yaml:"-"`
	ImagePrefix string `yaml:"-"`
}

func DefaultDambfile() *Dambfile {
	return &Dambfile{
		Dockerfile: "Dockerfile",
		Args: map[string]*Arg{
			"damb_tag":    {Cmd: []string{"git", "symbolic-ref", "--short", "HEAD"}},
			"damb_prefix": {Val: "local/"},
		},
		CacheFromTags: []string{"master"},
		BuildCmd:      []string{"docker", "build"},
	}
}

func (a *Arg) Get() (string, error) {
	if a == nil {
		return "", fmt.Errorf("argument is not set")
	}
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

func (cfg *Dambfile) Recompute() error {
	var err error
	cfg.ImagePrefix, err = cfg.Args["damb_prefix"].Get()
	if cfg.ImagePrefix == "" || err != nil {
		return fmt.Errorf("Unable to determine damb_prefix (%q)", err)
	}
	cfg.ImageTag, err = cfg.Args["damb_tag"].Get()
	if cfg.ImageTag == "" || err != nil {
		return fmt.Errorf("Unable to determine damb_tag (%q)", err)
	}
	return nil
}
