package main

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tssd",
		Short: "TSSD is a tool for playing around with threshold signatures schemes",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())
			return nil
		},
	}
	initRootCmd(rootCmd)
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		NewInitCmd(),
		NewKeygenSimulateCmd(),
		NewKeysignSimulateCmd(),
	)
}
