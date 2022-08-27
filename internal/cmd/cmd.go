package cmd

import "github.com/spf13/cobra"

func NewCli() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tm",
		Short: "The simplest task manager",
	}

	cmd.AddCommand(NewDoCommand())

	return cmd
}
