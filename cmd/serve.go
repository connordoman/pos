package cmd

import (
	"io"
	"log"
	"net/http"

	"github.com/connordoman/pos/internal/escpos"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
)

var ServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start the POS printer server",
	Long:  "Start the POS printer server and listen for print jobs",
	RunE:  runServeCommand,
}

func init() {

}

func runServeCommand(cmd *cobra.Command, args []string) error {
	r := chi.NewMux()

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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Post("/print", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("read error: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		text := string(bodyBytes)
		p.WriteString(text)

		w.Write([]byte("Print job received"))
	})

	return http.ListenAndServe(":42069", r)
}
