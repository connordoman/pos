package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/connordoman/pos/internal/escpos"
)

func handleCut(w http.ResponseWriter, r *http.Request) {
	defer lockPrinter()()

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

	if err := p.Cut(); err != nil {
		log.Printf("printer cut error: %v", err)
		http.Error(w, "Failed to cut paper", http.StatusInternalServerError)
		return
	}

	bytes, err := p.Flush()
	if err != nil {
		log.Printf("printer flush error: %v", err)
		http.Error(w, "Failed to send cut command to printer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Cut command sent to printer: %v", bytes)
}
