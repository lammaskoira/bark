/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/lammaskoira/bark/pkg/eval"
	"github.com/lammaskoira/bark/pkg/input"
	"github.com/lammaskoira/bark/pkg/parser"
)

var barkCmd = newBarkCmd()

func newBarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bark",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		SilenceUsage: true,
		RunE:         bark,
	}

	cmd.Flags().StringP("input", "i", "", "input path")
	cmd.Flags().StringP("trickset", "t", "", "trickset path")
	cmd.Flags().StringP("repodir", "r", "", "repo dir path")
	return cmd
}

// nolint:gochecknoinits // This is used by cobra
func init() {
	rootCmd.AddCommand(barkCmd)
}

func bark(cmd *cobra.Command, args []string) error {
	inputPath, iferr := notEmptyStringFlag(cmd, "input")
	if iferr != nil {
		return iferr
	}

	policyPath, pferr := notEmptyStringFlag(cmd, "trickset")
	if pferr != nil {
		return pferr
	}

	repoDirPath, rderr := notEmptyStringFlag(cmd, "repodir")
	if rderr != nil {
		return rderr
	}

	i, ierr := input.GetInputFromFile(inputPath)
	if ierr != nil {
		return ierr
	}

	pf, oerr := os.Open(policyPath)
	if oerr != nil {
		return fmt.Errorf("failed to open policy file: %w", oerr)
	}

	defer pf.Close()

	p := parser.NewParser()
	ts, perr := p.Parse(pf)
	if perr != nil {
		return fmt.Errorf("couldn't parse TrickSet file: %w", perr)
	}

	if verr := ts.ValidateRules(); verr != nil {
		return fmt.Errorf("invalid policy: %w", verr)
	}

	if chrerr := syscall.Chroot(repoDirPath); chrerr != nil {
		return fmt.Errorf("failed to chroot to repo dir: %w", chrerr)
	}

	if err := eval.EvaluateRules(cmd.Context(), i, ts.Rules); err != nil {
		return err
	}

	return nil
}
