package cmd

import (
	"log"

	"github.com/connordoman/pos/internal/escpos"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "pos",
	Short: "A CLI for handling POS printing commands",
	RunE:  runRootCommand,
}

func init() {
	RootCommand.AddCommand(
		ServeCommand,
	)
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	p, err := escpos.InitPrinter()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}()

	p.Init()
	p.SimpleLine()
	p.Feed(3)
	p.Log("Hello, world!")
	p.Feed(3)
	p.SimpleLine()
	p.FeedAndCut(3)
	p.Flush()
	return nil
}
