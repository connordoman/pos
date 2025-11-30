package server

import (
	"log"
	"net/http"

	"github.com/connordoman/pos/internal/supabase"
)

func handleGetJobs(w http.ResponseWriter, r *http.Request) {
	client, err := supabase.NewClient()
	if err != nil {
		log.Printf("error creating supabase client: %v", err)
		http.Error(w, "Failed to create supabase client", http.StatusInternalServerError)
		return
	}

	jobs, count, err := client.From("jobs").Select("*", "exact", false).Execute()
	if err != nil {
		log.Printf("error getting jobs: %v", err)
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jobs)
}
