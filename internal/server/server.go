package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/connordoman/pos/internal/escpos"
	"github.com/go-chi/chi"
)

var printMu sync.Mutex

type contextKey string

const printerKey contextKey = "printer"

type Server struct {
	r *chi.Mux

	Port string
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s.r)
}

func NewServer(port string) *Server {
	r := chi.NewMux()

	r.Get("/health", handleHealth)

	r.With(printerMiddleware).Post("/print", handlePrint)

	r.With(printerMiddleware).Post("/print/markdown", handlePrintMarkdown)

	r.With(printerMiddleware).Post("/cut", handleCut)

	return &Server{
		r:    r,
		Port: port,
	}
}

// printerMiddleware initializes the printer, locks it, and makes it available via context.
// It ensures the printer is properly closed and unlocked when the request completes.
func printerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Lock the printer mutex
		printMu.Lock()

		// Initialize the printer
		p, err := escpos.InitPrinter()
		if err != nil {
			printMu.Unlock()
			log.Printf("init printer error: %v", err)
			http.Error(w, "Printer not available", http.StatusServiceUnavailable)
			return
		}

		// Initialize the printer device
		if err := p.Init(); err != nil {
			printMu.Unlock()
			log.Printf("printer init error: %v", err)
			if closeErr := p.Close(); closeErr != nil {
				log.Printf("printer close error: %v", closeErr)
			}
			http.Error(w, "Failed to initialize printer", http.StatusInternalServerError)
			return
		}

		// Store printer in context
		ctx := context.WithValue(r.Context(), printerKey, p)
		r = r.WithContext(ctx)

		// Defer unlock and close
		defer func() {
			printMu.Unlock()
			if err := p.Close(); err != nil {
				log.Printf("printer close error: %v", err)
			}
		}()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// getPrinter retrieves the printer from the request context.
// It panics if the printer is not found (which should never happen if middleware is used correctly).
func getPrinter(r *http.Request) *escpos.Printer {
	p, ok := r.Context().Value(printerKey).(*escpos.Printer)
	if !ok {
		panic("printer not found in context - ensure printerMiddleware is used")
	}
	return p
}

func lockPrinter() func() {
	printMu.Lock()

	return func() {
		printMu.Unlock()
	}
}
