package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stiletto/damb/repo"
)

func (damb *Damb) ExpandAlias(alias string) []string {
	result := make([]string, 0, 1)
	if expanded, ok := damb.cfg.Aliases[alias]; ok {
		result = append(result, expanded...)
	} else {
		result = append(result, alias)
	}
	return result
}

func (damb *Damb) NormalizeTargets(args []string) []string {
	targets := make([]string, 0, len(args))
	for _, arg := range args {
		for _, tgt := range damb.ExpandAlias(arg) {
			spec, err := damb.cfg.ParseLocalTarget(tgt)
			if err != nil {
				damb.Fatalf(1, " %s: %s", arg, err)
			} else {
				damb.Printf(damb.au.Green(" %s: %s"), arg, spec.Image)
				targets = append(targets, spec.Image)
			}
		}
	}
	return targets
}

func (damb *Damb) Resolve(cmd *cobra.Command, args []string) {
	damb.Init(cmd)
	targets := damb.NormalizeTargets(args)
	deptype, err := cmd.Flags().GetString("type")
	if err != nil {
		damb.Fatalf(2, "%s", err)
	}
	depall := false
	var depforeign bool
	switch deptype {
	case "all":
		depall = true
	case "internal":
		depforeign = false
	case "external":
		depforeign = true
	default:
		damb.Fatalf(1, "Unsupported target type: %q", deptype)
	}

	err = damb.repo.Resolve(targets, func(t *repo.Target) error {
		if depall || depforeign == t.IsForeign() {
			fmt.Println(t.Image)
		}
		return nil
	})
	if err != nil {
		damb.Fatalf(1, "%s", err)
	}
}
