package cmd

import (
	"io"
	"log"
	"net/http"
    "strings"
    "sync"

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

// printMu serializes access to the USB printer to avoid concurrent writes
// from multiple HTTP requests.
var printMu sync.Mutex

func runServeCommand(cmd *cobra.Command, args []string) error {
	r := chi.NewMux()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Post("/print", func(w http.ResponseWriter, r *http.Request) {
		// Ensure only one request talks to the device at a time
		printMu.Lock()
		defer printMu.Unlock()

		p, err := escpos.InitPrinter()
		if err != nil {
			log.Printf("init printer error: %v", err)
			http.Error(w, "Printer not available", http.StatusServiceUnavailable)
			return
		}
		defer func() {
			if err := p.Close(); err != nil {
				log.Printf("printer close error: %v", err)
			}
		}()

		if err := p.Init(); err != nil {
			log.Printf("printer init error: %v", err)
			http.Error(w, "Failed to initialize printer", http.StatusInternalServerError)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("read body error: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		text := string(bodyBytes)
		if len(text) == 0 {
			http.Error(w, "Empty print job", http.StatusBadRequest)
			return
		}
		// Ensure at least one newline so output becomes visible on paper
		if !strings.HasSuffix(text, "\n") {
			text += "\n"
		}

		p.WriteString(text)
		// Feed a couple lines to push content out of the head area
		p.Feed(2)

		if _, err := p.Flush(); err != nil {
			log.Printf("flush error: %v", err)
			http.Error(w, "Failed to send data to printer", http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Print job queued to device"))
	})

	return http.ListenAndServe(":42069", r)
}
