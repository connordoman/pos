package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/connordoman/pos/internal/escpos"
)

func handlePrint(w http.ResponseWriter, r *http.Request) {
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
	p.FeedAndCut(10)

	if _, err := p.Flush(); err != nil {
		log.Printf("flush error: %v", err)
		http.Error(w, "Failed to send data to printer", http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Print job queued to device"))

}

func handlePrintMarkdown(w http.ResponseWriter, r *http.Request) {
	defer lockPrinter()()

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

	if err := p.ParseMarkdown(text); err != nil {
		log.Printf("printer parse markdown error: %v", err)
		http.Error(w, "Failed to parse markdown", http.StatusBadRequest)
		return
	}

	if err := p.FeedAndCut(5); err != nil {
		log.Printf("printer feed and cut error: %v", err)
		http.Error(w, "Failed to feed and cut paper", http.StatusInternalServerError)
		return
	}

	b, err := p.Flush()
	if err != nil {
		log.Printf("printer flush error: %v", err)
		http.Error(w, "Failed to send data to printer", http.StatusBadGateway)
		return
	}

	msg := fmt.Sprintf("Print job queued to device (%d bytes)", b)

	log.Println(msg)

	w.Write([]byte(msg))
}
