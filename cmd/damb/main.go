package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/stiletto/damb/dambfile"
	"github.com/stiletto/damb/dockerfile"
)

type Target struct {
	Image        string
	Directory    string
	Dockerfile   string
	Stage        string
	Dependencies []string
	Seen         bool
	Built        bool
	Args         []string
}

func LoadTargets(cfg *dambfile.Dambfile, directory string) ([]*Target, error) {
	fmt.Printf("Discovering targets in %s\n", directory)
	dir, err := os.Open(filepath.Join(cfg.DirRoot, directory))
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fis, err := dir.Readdir(4096)
	if err != nil {
		return nil, err
	}
	targets := make([]*Target, 0)
	for _, fi := range fis {
		name := fi.Name()
		if name == cfg.Dockerfile || strings.HasPrefix(name, cfg.Dockerfile+".") {
			df, err := dockerfile.Load(cfg.Args, filepath.Join(cfg.DirRoot, directory, name))
			if err != nil {
				return nil, err
			}
			stagePrefix := name[len(cfg.Dockerfile):]
			knownStages := make(map[string]*Target)
			for _, stage := range df.Stages {
				stageName := stage.Name
				if stageName != "" {
					stageName = "." + stageName
				}
				stageName = stagePrefix + stageName
				target := &Target{
					Image:        path.Join(cfg.ImagePrefix, directory+stageName) + ":" + cfg.Args["damb_ref"].Val,
					Directory:    directory,
					Dockerfile:   name,
					Stage:        stage.Name,
					Args:         make([]string, 0),
					Dependencies: make([]string, 0),
				}
				for dep := range stage.Dependencies {
					if depTarget := knownStages[dep]; depTarget != nil {
						target.Dependencies = append(target.Dependencies, depTarget.Image)
					} else {
						target.Dependencies = append(target.Dependencies, dep)
					}
				}
				for arg := range stage.Args {
					target.Args = append(target.Args, arg)
				}
				fmt.Printf(" - %s\n", target.Image)
				targets = append(targets, target)
				knownStages[stage.Name] = target
			}
		}
	}
	return targets, nil
}

func ImageToDirectory(cfg *dambfile.Dambfile, image string) (string, string, bool) {
	if !strings.HasPrefix(image, cfg.ImagePrefix) {
		return "", "", false
	}
	image = image[len(cfg.ImagePrefix):]
	if !strings.HasSuffix(image, ":"+cfg.RefName) {
		return "", "", false
	}
	dirPlusStage := image[:len(image)-len(cfg.RefName)-1]
	delimiter := strings.Index(dirPlusStage, ".")
	if delimiter == -1 {
		delimiter = len(dirPlusStage)
	}
	return dirPlusStage[:delimiter], dirPlusStage[delimiter:], true
}

func (c *Cache) Build(image string) error {
	target, err := c.Get(image)
	if err != nil {
		if err == ErrForeign {
			return nil
		}
		return err
	}
	if target.Built {
		return nil
	}
	if target.Seen {
		return fmt.Errorf("Found a loop somewhere around %q", target.Image)
	}
	target.Seen = true
	for _, dep := range target.Dependencies {
		err := c.Build(dep)
		if err != nil {
			return err
		}
	}
	if len(c.cfg.BuildCmd) < 1 {
		return fmt.Errorf("Build command is empty")
	}
	buildCmd := make([]string, 0, 8)
	buildCmd = append(buildCmd, c.cfg.BuildCmd...)
	buildCmd = append(buildCmd, "-t", target.Image)
	buildCmd = append(buildCmd, "-f", target.Dockerfile)
	for _, k := range target.Args {
		arg := c.cfg.Args[k]
		if arg != nil {
			argValue, err := arg.Get()
			if err != nil {
				return fmt.Errorf("%q: %s", k, err)
			}
			buildCmd = append(buildCmd, "--build-arg", k+"="+argValue)
		} else {
			fmt.Printf("Argument %q is probably set to default value\n", k)
		}
	}
	if target.Stage != "" {
		buildCmd = append(buildCmd, "--target", target.Stage)
	}
	buildCmd = append(buildCmd, ".")
	fmt.Printf("buildCmd: %#v\n", buildCmd)
	cmd := exec.Command(buildCmd[0], buildCmd[1:]...)
	cmd.Dir = target.Directory
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	target.Built = true
	return nil
}

func main() {
	cfg, err := dambfile.FindAndLoad(".")
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}
	fmt.Printf("%#v\n", cfg)
	err = os.Chdir(cfg.DirRoot)
	if err != nil {
		log.Fatalf("Failed to chdir to repo root: %s", err)
	}

	if len(os.Args) > 1 {
		c := &Cache{
			cfg:         cfg,
			targets:     make(map[string]*Target),
			directories: make(map[string]struct{}),
		}
		for _, arg := range os.Args[1:] {
			if !strings.HasPrefix(arg, cfg.ImagePrefix) {
				arg = cfg.ImagePrefix + arg
			}
			if !strings.Contains(arg, ":") {
				arg = arg + ":" + cfg.RefName
			}
			tgt, stg, ok := ImageToDirectory(cfg, arg)
			if !ok { // this may happen if target has different tag
				log.Fatalf("Refusing to build %q", arg)
			}
			tgt = strings.TrimLeft(filepath.Clean(tgt), "/")
			tgt = cfg.ImagePrefix + tgt + stg + ":" + cfg.RefName
			fmt.Printf("Building %s\n", tgt)

			err := c.Build(tgt)
			if err != nil {
				log.Fatalf("Failed to build %q: %s", arg, err)
			}
		}
	}
	/*		targets, err := LoadTargets(cfg, os.Args[1])
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			for _, target := range targets {
				fmt.Printf("Target %s:\n", target.Image)
				fmt.Printf("  Directory:  %s\n", target.Directory)
				fmt.Printf("  Dockerfile: %s\n", target.Dockerfile)
				fmt.Printf("  Stage:      %s\n", target.Stage)
				fmt.Printf("  Dependencies:\n")
				for _, dep := range target.Dependencies {
					deptype := "local"
					if _, ok := ImageToDirectory(cfg, dep); !ok {
						deptype = "foreign"
					}
					fmt.Printf("    - %s (%s)\n", dep, deptype)
				}
			}
		}*/
}
