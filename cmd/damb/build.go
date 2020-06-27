package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/stiletto/damb/repo"
)

func (damb *Damb) Build(cmd *cobra.Command, args []string) {
	damb.Init(cmd)
	targets := damb.NormalizeTargets(args)
	nocap, err := cmd.Flags().GetBool("no-capture")
	if err != nil {
		damb.Fatalf(2, "%s", err)
	}

	err = damb.repo.Resolve(targets, func(t *repo.Target) error {
		if t.IsForeign() {
			return nil
		}
		if !nocap {
			fmt.Fprint(os.Stderr, damb.au.Sprintf(damb.au.BrightMagenta("[building] %s"), t.Image))
		}
		var buildOut bytes.Buffer
		var ctxOut bytes.Buffer
		err := func() error {
			if len(damb.cfg.BuildCmd) < 1 {
				return fmt.Errorf("Build command is empty")
			}
			buildCmd := make([]string, 0, 8)
			buildCmd = append(buildCmd, damb.cfg.BuildCmd...)
			buildCmd = append(buildCmd, "-t", t.Image)
			for k := range t.Args {
				arg := damb.cfg.Args[k]
				if arg != nil {
					argValue, err := arg.Get()
					if err != nil {
						return fmt.Errorf("%q: %s", k, err)
					}
					buildCmd = append(buildCmd, "--build-arg", k+"="+argValue)
				} /*else {
					fmt.Printf("Argument %q is probably set to default value\n", k)
				}*/
			}
			/* if target.Stage != "" {
				buildCmd = append(buildCmd, "--target", target.Stage)
			}*/
			damb.Debugf("")
			ctxCmd := make([]string, 0, 8)
			if t.Options.Context == "" {
				buildCmd = append(buildCmd, "-f", t.Dockerfile)
				buildCmd = append(buildCmd, ".")
			} else {
				buildCmd = append(buildCmd, "-f", filepath.Join(t.Directory, t.Dockerfile))
				buildCmd = append(buildCmd, "-")
				ctxCmd = []string{"sh", "-c", t.Options.Context}
				damb.Debugf("           ctxCmd: %#v", ctxCmd)
			}
			damb.Debugf("           buildCmd: %#v", buildCmd)
			buildProc := exec.Command(buildCmd[0], buildCmd[1:]...)
			if nocap {
				buildProc.Stdout = os.Stdout
				buildProc.Stderr = os.Stderr
			} else {
				buildProc.Stdout = &buildOut
				buildProc.Stderr = &buildOut
			}
			var ctxProc *exec.Cmd
			if t.Options.Context == "" {
				buildProc.Dir = filepath.Join(damb.cfg.DirRoot, t.Directory)
				buildProc.Stdin = nil
			} else {
				buildProc.Dir = filepath.Join(damb.cfg.DirRoot)
				ctxProc = exec.Command(ctxCmd[0], ctxCmd[1:]...)
				ctxProc.Stdin = nil
				if nocap {
					ctxProc.Stderr = os.Stderr
				} else {
					ctxProc.Stderr = &ctxOut
				}
				buildProc.Stdin, err = ctxProc.StdoutPipe()
				if err != nil {
					return fmt.Errorf("ctxCmd StdoutPipe: %w", err)
				}
				if err = ctxProc.Start(); err != nil {
					return fmt.Errorf("%#v: %w", ctxCmd, err)
				}
			}
			err := buildProc.Run()
			var ctxErr error
			if ctxProc != nil {
				ctxErr = ctxProc.Wait()
			}
			if err == nil {
				if ctxErr != nil {
					err = fmt.Errorf("%#v: %w", ctxCmd, ctxErr)
				}
			} else {
				err = fmt.Errorf("%#v: %w", buildCmd, err)
			}
			if !nocap {
				fmt.Fprint(os.Stderr, "\r")
			}
			return err
		}()
		resultColor := damb.au.BrightGreen
		if err != nil {
			resultColor = damb.au.BrightRed
		}

		fmt.Fprintln(os.Stderr, damb.au.Sprintf(resultColor("[%8s] %s"), time.Since(damb.start), t.Image))
		if err != nil && !nocap {
			damb.Printf(damb.au.BrightRed("Failed to build image for %q:\n%s%s"), t.Image, damb.au.Red(ctxOut.String()), damb.au.Red(buildOut.String()))
		}
		return err
	})
	if err != nil {
		damb.Fatalf(1, "%s", err)
	}
}
