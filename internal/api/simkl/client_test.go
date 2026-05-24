package simkl

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMediaDetailsEmptyArrayReturnsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := NewClient(NewClientOptions{
		ClientID:   "test-client",
		APIBaseURL: server.URL,
		HTTPClient: server.Client(),
	})

	_, err := client.MediaDetails(context.Background(), MediaTypeMovies, "123", "full")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for empty media-details array, got %v", err)
	}
}

func TestMediaEpisodesAllowsEmptyArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := NewClient(NewClientOptions{
		ClientID:   "test-client",
		APIBaseURL: server.URL,
		HTTPClient: server.Client(),
	})

	episodes, err := client.MediaEpisodes(context.Background(), MediaTypeShows, "123", "full")
	if err != nil {
		t.Fatalf("expected empty episode arrays to decode without error, got %v", err)
	}
	if len(episodes) != 0 {
		t.Fatalf("expected no episodes, got %d", len(episodes))
	}
}
