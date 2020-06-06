package dambfile

import (
	"os"
	"path/filepath"

	"github.com/stiletto/damb/utils"
	"gopkg.in/yaml.v3"
)

var DambfileNames = []string{"Damb.yml", ".damb.yml"}

func Find(start string) (string, error) {
	return utils.FindUp(start, DambfileNames)
}

func Load(fname string) (*Dambfile, error) {
	absfname, err := filepath.Abs(fname)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(absfname)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	cfg := DefaultDambfile()
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}
	cfg.DirRoot = filepath.Dir(absfname)
	if cfg.Args == nil {
		cfg.Args = make(map[string]*Arg)
	}
	cfg.RefName, err = utils.CommandOutput(cfg.RefNameCmd...)
	if err != nil {
		return nil, err
	}
	cfg.Args["damb_root"] = &Arg{Val: cfg.ImagePrefix}
	cfg.Args["damb_ref"] = &Arg{Val: cfg.RefName}
	return cfg, nil
}

func FindAndLoad(start string) (*Dambfile, error) {
	fname, err := Find(start)
	if err != nil {
		return nil, err
	}
	return Load(fname)
}
