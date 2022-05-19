/*
Copyright © 2022 Lammaskoira authors

*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bark",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
