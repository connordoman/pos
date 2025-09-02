package cmd

import (
	"fmt"
	"log"

	"github.com/connordoman/pos/internal/escpos"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var RSSCommand = &cobra.Command{
	Use:   "rss",
	Short: "Run RSS feed commands",
	Long:  "long description",
	RunE:  runRSSCommand,
}

func init() {

}

func runRSSCommand(cmd *cobra.Command, args []string) error {

	p, err := escpos.InitPrinter()
	if err != nil {
		return err
	}
	defer func() {
		if err := p.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}()

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://news.ycombinator.com/rss")
	fmt.Println(feed.Title)

	fmt.Println()

	for _, item := range feed.Items {
		fmt.Printf("Title: %s\nLink: %s\n\n", item.Title, item.Link)
	}

	return nil
}
