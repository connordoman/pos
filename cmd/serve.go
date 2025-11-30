package cmd

import (
	"log"
	"strconv"

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

func runServeCommand(cmd *cobra.Command, args []string) error {
	portFlag, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	port := strconv.Itoa(portFlag)

	s := server.NewServer(port)

	log.Printf("Starting server on :%s", port)

	return s.Start()
}
