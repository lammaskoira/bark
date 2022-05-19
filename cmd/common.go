package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func notEmptyStringFlag(cmd *cobra.Command, name string) (string, error) {
	flag, ferr := cmd.Flags().GetString(name)
	if ferr != nil {
		return "", fmt.Errorf("flag %s not set", name)
	}

	if flag == "" {
		return "", fmt.Errorf("flag %s is empty", name)
	}

	return flag, nil
}
