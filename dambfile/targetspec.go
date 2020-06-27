package dambfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TargetSpec struct {
	Image      string
	Directory  string
	Dockerfile string
}

func (ts *TargetSpec) IsForeign() bool {
	return ts.Directory == ""
}

func parseDirStage(dirPlusStage string) (directory string, stage string) {
	if delimiter := strings.Index(dirPlusStage, "."); delimiter != -1 {
		directory = dirPlusStage[:delimiter]
		stage = dirPlusStage[delimiter+1:]
	} else {
		directory = dirPlusStage
		stage = ""
	}
	directory = filepath.FromSlash(directory)
	return directory, stage
}

func (cfg *Dambfile) ParseTarget(image string) (target TargetSpec) {
	target = TargetSpec{
		Image:     image,
		Directory: "", // empty Direcyory means "foreign"
	}
	if !strings.HasPrefix(image, cfg.ImagePrefix) {
		return target // "", "", fmt.Errorf("Image doesn't start with %q prefix: %q", c.cfg.ImagePrefix, image)
	}
	image = image[len(cfg.ImagePrefix):]
	if !strings.HasSuffix(image, ":"+cfg.ImageTag) {
		return target // "", "", fmt.Errorf("Image doesn't end with :%q suffix: %q", c.cfg.RefName, image)
	}
	dirPlusStage := image[:len(image)-len(cfg.ImageTag)-1]
	var stage string
	target.Directory, stage = parseDirStage(dirPlusStage)
	if stage == "" {
		target.Dockerfile = cfg.Dockerfile
	} else {
		target.Dockerfile = cfg.Dockerfile + "." + stage
	}
	return target
}

func (cfg *Dambfile) ParseLocalTarget(targetName string) (TargetSpec, error) {
	target := cfg.ParseTarget(targetName)
	if !target.IsForeign() {
		return target, nil
	}

	directory, stage := parseDirStage(targetName)
	directoryAbs := filepath.Join(cfg.DirRoot, directory)
	fi, err := os.Stat(directoryAbs)
	if err != nil {
		return target, err
	}
	if !fi.IsDir() {
		return target, fmt.Errorf("%q is not a directory", directoryAbs)
	}
	dockerfile := cfg.Dockerfile
	if stage != "" {
		dockerfile += "." + stage
	}
	return TargetSpec{
		Image:      cfg.ImagePrefix + targetName + ":" + cfg.ImageTag,
		Directory:  directory,
		Dockerfile: dockerfile,
	}, nil
}
