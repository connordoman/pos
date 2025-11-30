package cmd

import (
	"log"
	"sync"

	"github.com/connordoman/pos/internal/server"
	"github.com/spf13/cobra"
)

var ServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start the POS printer server",
	Long:  "Start the POS printer server and listen for print jobs",
	RunE:  runServeCommand,
}

func init() {
	ServeCommand.Flags().IntP("port", "p", 42069, "The port to listen on")
}

// printMu serializes access to the USB printer to avoid concurrent writes
// from multiple HTTP requests.
var printMu sync.Mutex

func runServeCommand(cmd *cobra.Command, args []string) error {
	portFlag, err := cmd.Flags().GetString("port")
	if err != nil {
		return err
	}

	s := server.NewServer(portFlag)

	log.Printf("Starting server on :%s", portFlag)

	return s.Start()
}
