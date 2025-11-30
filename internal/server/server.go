package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var printMu sync.Mutex

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

	r.Post("/print", handlePrint)

	r.Post("/cut", handleCut)

	return &Server{
		r:    r,
		Port: port,
	}
}

func lockPrinter() func() {
	printMu.Lock()

	return func() {
		printMu.Unlock()
	}
}
