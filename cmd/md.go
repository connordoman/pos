package cmd

import (
	"github.com/connordoman/pos/internal/escpos/md"
	"github.com/spf13/cobra"
)

var MdCommand = &cobra.Command{
	Use:   "md [<file>]",
	Short: "Compile Markdown",
	Long:  "Compile Markdown to ESC/POS commands",
	RunE:  runMdCommand,
	Args:  cobra.MaximumNArgs(1),
}

func init() {

}

func runMdCommand(cmd *cobra.Command, args []string) error {
	interpreter := md.NewInterpreter()

	if len(args) == 0 {
		return interpreter.RunPrompt()
	}

	file := args[0]
	return interpreter.RunFile(file)
}
