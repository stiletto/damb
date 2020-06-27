package main

import (
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func main() {
	damb := &Damb{
		au: aurora.NewAurora(true),
	}
	dambCmd := &cobra.Command{
		Use:   "damb",
		Short: "Docker As Monorepo Build-system",
	}
	dambCmd.PersistentFlags().StringToStringP("arg", "a", make(map[string]string), "Override build arguments")
	dambCmd.PersistentFlags().Bool("debug", false, "Debug mode")
	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve and print target dependencies",
		Run:   damb.Resolve,
	}
	resolveCmd.Flags().StringP("type", "t", "internal", "Type of targets to list: external, internal or all")
	dambCmd.AddCommand(resolveCmd)
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build target images",
		Run:   damb.Build,
	}
	buildCmd.Flags().Bool("no-capture", false, "Don't capture docker output")
	buildCmd.Flags().String("build-cmd", "", "Override build_cmd")
	buildCmd.Flags().StringSlice("cache-tag", []string{}, "Consider these tags as cache sources")
	dambCmd.AddCommand(buildCmd)
	if err := dambCmd.Execute(); err != nil {
		damb.Fatalf(3, "%s", err)
	}
}
