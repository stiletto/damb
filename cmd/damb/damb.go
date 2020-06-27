package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/stiletto/damb/dambfile"
	"github.com/stiletto/damb/repo"
)

type Damb struct {
	initOnce sync.Once
	cfg      *dambfile.Dambfile
	debug    bool
	repo     *repo.Cache
	au       aurora.Aurora
	start    time.Time
}

func (damb *Damb) Printf(format interface{}, args ...interface{}) {
	fmt.Fprintln(os.Stderr, damb.au.Sprintf(format, args...))
}

func (damb *Damb) Fatalf(code int, format interface{}, args ...interface{}) {
	damb.Printf(damb.au.BrightRed(format), args...)
	os.Exit(code)
}

func (damb *Damb) Debugf(format interface{}, args ...interface{}) {
	if damb.debug {
		damb.Printf(damb.au.BrightBlue(format), args...)
	}
}

func (damb *Damb) Init(cmd *cobra.Command) {
	damb.initOnce.Do(func() {
		damb.start = time.Now()
		var err error
		damb.debug, _ = cmd.Flags().GetBool("debug")
		damb.Debugf("Debug mode is on")
		damb.cfg, err = dambfile.FindAndLoad(".")
		if err != nil {
			damb.Fatalf(1, "Failed to load configuration: %s", err)
		}
		if args, err := cmd.Flags().GetStringToString("arg"); err != nil {
			damb.Fatalf(1, "Failed to parse --arg: %s", err)
		} else {
			for k, v := range args {
				damb.cfg.Args[k] = &dambfile.Arg{Val: v}
			}
		}
		if err := damb.cfg.Recompute(); err != nil {
			damb.Fatalf(1, "Configuration failure: %s", err)
		}
		damb.Debugf("Configuration: %#v", damb.cfg)
		damb.Printf(damb.au.BrightYellow(fmt.Sprintf("Starting from %s", damb.cfg.DirRoot)))
		err = os.Chdir(damb.cfg.DirRoot)
		if err != nil {
			damb.Fatalf(1, "Failed to chdir to repo root: %s", err)
		}
		damb.repo = repo.New(damb.cfg)
	})
}
