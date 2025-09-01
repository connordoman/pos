package cmd

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:   "pos",
	Short: "A CLI for handling POS printing commands",
	RunE:  runRootCommand,
}

func init() {

}

func runRootCommand(cmd *cobra.Command, args []string) error {
	return nil
}
