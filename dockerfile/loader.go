package dockerfile

import (
	"fmt"
	"io"
	"os"
	"strconv"

	//	"github.com/openshift/imagebuilder/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/shell"

	"github.com/stiletto/damb/dambfile"
	"github.com/stiletto/damb/utils"
)

func Load(args map[string]*dambfile.Arg, fname string) (*Dockerfile, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	df, err := parse(args, f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fname, err)
	}
	return df, nil
}

func parse(args map[string]*dambfile.Arg, f io.Reader) (*Dockerfile, error) {
	res, err := parser.Parse(f)
	if err != nil {
		return nil, err
	}
	shlex := shell.NewLex('\\')
	stages, metaArgs, err := instructions.Parse(res.AST)
	if err != nil {
		return nil, err
	}
	knownMetaArgs := make(map[string]string)
	argExpander := func(word string) (string, error) {
		return shlex.ProcessWordWithMap(word, knownMetaArgs)
	}
	for _, arg := range metaArgs {
		err := arg.Expand(argExpander)
		if err != nil {
			return nil, err
		}
		if defArg, ok := args[arg.Key]; ok {
			knownMetaArgs[arg.Key], err = defArg.Get()
			if err != nil {
				return nil, err
			}
		} else {
			if arg.Value == nil {
				return nil, fmt.Errorf("Dockerfile uses undefined argument %q", arg.Key)
			} else {
				knownMetaArgs[arg.Key] = *arg.Value
			}
		}
	}
	df := &Dockerfile{
		Stages:   make([]Stage, 0, len(stages)),
		MetaArgs: knownMetaArgs,
	}
	for i, stage := range stages {
		expandedBase, err := argExpander(stage.BaseName)
		if err != nil {
			return nil, err
		}
		if stage.Name == "" && i < len(stages)-2 {
			return nil, fmt.Errorf("Only last stage is allowed to be unnamed (FROM %s)", stage.BaseName)
		}
		curStage := Stage{
			Name:         stage.Name,
			Dependencies: map[string]struct{}{expandedBase: {}},
			Args:         make(map[string]struct{}),
		}
		for _, arg := range metaArgs {
			curStage.Args[arg.Key] = struct{}{}
		}
		for _, cmd := range stage.Commands {
			if expandable, ok := cmd.(instructions.SupportsSingleWordExpansion); ok {
				err := expandable.Expand(argExpander)
				if err != nil {
					return nil, err
				}
			}
			switch cmd := cmd.(type) {
			case *instructions.CopyCommand:
				if cmd.From != "" {
					if _, err := strconv.Atoi(cmd.From); err == nil {
						return nil, fmt.Errorf("Don't reference stages by their number (%s)", cmd.String())
					}
					curStage.Dependencies[cmd.From] = struct{}{}
				}
			case *instructions.ArgCommand:
				if cmd.Value == nil && args[cmd.Key] == nil {
					return nil, fmt.Errorf("Stage %q (FROM %s) uses undefined argument: %q", stage.Name, stage.BaseName, cmd.Key)
				}
				curStage.Args[cmd.Key] = struct{}{}
			}
		}
		delete(curStage.Dependencies, "scratch")
		df.Stages = append(df.Stages, curStage)
	}
	return df, err
}

func FindAndLoad(args map[string]*dambfile.Arg, start, name string) (*Dockerfile, error) {
	fname, err := utils.FindUp(start, []string{name})
	if err != nil {
		return nil, err
	}
	return Load(args, fname)
}
