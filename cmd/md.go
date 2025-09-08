package cmd

import (
	"fmt"
	"os"

	"github.com/connordoman/pos/internal/escpos"
	"github.com/spf13/cobra"
)

var MdCommand = &cobra.Command{
	Use:   "md [file]",
	Short: "Compile Markdown to ESC/POS bytes and print to stdout",
	RunE:  runMdCommand,
	Args:  cobra.MaximumNArgs(1),
}

func init() {}

func runMdCommand(cmd *cobra.Command, args []string) error {
	path := "templates/test-markdown.md"
	if len(args) > 0 && args[0] != "" {
		path = args[0]
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Build ESC/POS bytes into buffer without talking to USB
	p := &escpos.Printer{}
	if err := p.SelectFont('A'); err != nil {
		return err
	}
	if err := p.ParseMarkdown(string(data)); err != nil {
		return err
	}
	// Preview to stdout
	fmt.Print(string(p.Buff()))
	return nil
}
