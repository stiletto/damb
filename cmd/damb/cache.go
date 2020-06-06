package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/stiletto/damb/dambfile"
)

type Cache struct {
	cfg         *dambfile.Dambfile
	targets     map[string]*Target
	directories map[string]struct{}
}

var ErrForeign = errors.New("Foreign target")

func (c *Cache) Get(image string) (*Target, error) {
	target := c.targets[image]
	if target != nil {
		return target, nil
	}
	dir, _, ok := ImageToDirectory(c.cfg, image)
	if !ok {
		return nil, ErrForeign
	}
	if filepath.Clean(dir) != dir {
		return nil, fmt.Errorf("Malformed image name: %s", image)
	}

	if _, ok := c.directories[dir]; !ok {
		targets, err := LoadTargets(c.cfg, dir)
		if err != nil {
			return nil, err
		}
		for _, t := range targets {
			if ot := c.targets[t.Image]; ot != nil {
				return nil, fmt.Errorf("Multiple targets with name %q: %#v and %#v\n", t.Image, ot, t)
			}
		}
		c.directories[dir] = struct{}{}
		for _, t := range targets {
			c.targets[t.Image] = t
		}
	}

	target = c.targets[image]
	if target != nil {
		return target, nil
	}

	return nil, fmt.Errorf("Target %q not found", image)
}
