package supabase

import (
	"errors"
	"os"

	"github.com/supabase-community/supabase-go"
)

func init() {

}

func NewClient() (*supabase.Client, error) {
	sbUrl := os.Getenv("SUPABASE_URL")
	sbKey := os.Getenv("SUPABASE_KEY")
	if sbUrl == "" || sbKey == "" {
		return nil, errors.New("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	return supabase.NewClient(sbUrl, sbKey, nil)
}
