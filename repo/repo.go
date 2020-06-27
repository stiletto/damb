package repo

import (
	"fmt"
	"path/filepath"

	"github.com/stiletto/damb/dambfile"
)

type Cache struct {
	cfg     *dambfile.Dambfile
	targets map[string]*Target
}

func New(cfg *dambfile.Dambfile) *Cache {
	return &Cache{
		cfg:     cfg,
		targets: make(map[string]*Target),
	}
}

func (c *Cache) Get(image string) (*Target, error) {
	target := c.targets[image]
	if target != nil {
		return target, nil
	}
	spec := c.cfg.ParseTarget(image)
	if spec.IsForeign() {
		target = &Target{TargetSpec: spec, Dependencies: make(map[string]struct{}), Args: make(map[string]struct{})}
	} else {
		if filepath.Clean(spec.Directory) != spec.Directory {
			return nil, fmt.Errorf("Malformed image name: %s", image)
		}
		var err error
		target, err = LoadTarget(c.cfg.DirRoot, spec, c.cfg.Args)
		if err != nil {
			return nil, err
		}
	}
	c.targets[image] = target
	return target, nil
}

type ResolveCallback func(*Target) error

func (c *Cache) resolveInternal(image string, out ResolveCallback, entered, exited map[*Target]struct{}) error {
	target, err := c.Get(image)
	if err != nil {
		if err == ErrForeign {
			return nil
		}
		return fmt.Errorf("%s: %w", image, err)
	}
	if _, ok := exited[target]; ok {
		return nil
	}
	if _, ok := entered[target]; ok {
		return fmt.Errorf("Found a loop somewhere around %q", target.Image)
	}
	entered[target] = struct{}{}
	for dep := range target.Dependencies {
		err := c.resolveInternal(dep, out, entered, exited)
		if err != nil {
			return err
		}
	}
	err = out(target)
	if err != nil {
		err = fmt.Errorf("%s: %w", image, err)
	}
	exited[target] = struct{}{}
	return err
}

func (c *Cache) Resolve(images []string, out ResolveCallback) error {
	entered := make(map[*Target]struct{})
	exited := make(map[*Target]struct{})
	for _, image := range images {
		err := c.resolveInternal(image, out, entered, exited)
		if err != nil {
			return err
		}
	}
	return nil
}
