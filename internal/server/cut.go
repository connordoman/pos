package server

import (
	"fmt"
	"log"
	"net/http"
)

func handleCut(w http.ResponseWriter, r *http.Request) {
	p := getPrinter(r)

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
