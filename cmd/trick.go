/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lammaskoira/bark/pkg/parser"
	"github.com/lammaskoira/bark/pkg/runner"
)

// trickCmd represents the trick command.
var trickCmd = newTrickCmd()

func newTrickCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trick",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		SilenceUsage: true,
		RunE:         trick,
	}

	cmd.Flags().StringP("trickset", "t", "", "trickset file to evaluate")
	return cmd
}

// nolint:gochecknoinits // This is used by cobra.
func init() {
	rootCmd.AddCommand(trickCmd)
}

func getTrickSetPathFromParam(cmd *cobra.Command) (string, error) {
	tsfpath, iferr := cmd.Flags().GetString("trickset")
	if iferr != nil {
		return "", iferr
	}

	if tsfpath == "" {
		return "", fmt.Errorf("no TrickSet file path specified")
	}

	return tsfpath, nil
}

func trick(cmd *cobra.Command, args []string) error {
	tsfpath, iferr := getTrickSetPathFromParam(cmd)
	if iferr != nil {
		return iferr
	}

	tsf, oerr := os.Open(tsfpath)
	if oerr != nil {
		return fmt.Errorf("could not open TrickSet file: %w", oerr)
	}
	defer tsf.Close()

	p := parser.NewParser()
	ts, perr := p.Parse(tsf)
	if perr != nil {
		return fmt.Errorf("couldn't parse TrickSet file: %w", perr)
	}

	return runner.Run(cmd.Context(), ts)
}
