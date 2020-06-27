package repo

import (
	"errors"
	"path/filepath"

	"github.com/stiletto/damb/dambfile"
	"github.com/stiletto/damb/dockerfile"
)

type Target struct {
	dambfile.TargetSpec
	// Stage        string
	Options      DambDockerfileOptions
	Dependencies map[string]struct{}
	Args         map[string]struct{}
}

var ErrForeign = errors.New("Foreign target")

func LoadTarget(root string, spec dambfile.TargetSpec, args map[string]*dambfile.Arg) (*Target, error) {
	if spec.IsForeign() {
		return nil, ErrForeign
	}
	target := &Target{
		TargetSpec:   spec,
		Args:         make(map[string]struct{}),
		Dependencies: make(map[string]struct{}),
	}
	dockerfilePath := filepath.Join(root, spec.Directory, spec.Dockerfile)
	df, err := dockerfile.Load(args, dockerfilePath)
	if err != nil {
		return nil, err
	}
	knownStages := make(map[string]struct{})
	for _, stage := range df.Stages {
		for dep := range stage.Dependencies {
			if _, ok := knownStages[dep]; !ok {
				target.Dependencies[dep] = struct{}{}
			}
		}
		for arg := range stage.Args {
			target.Args[arg] = struct{}{}
		}
		knownStages[stage.Name] = struct{}{}
	}
	target.Options, err = LoadDambOptions(dockerfilePath)
	if err != nil {
		return nil, err
	}
	return target, nil
}
