package alert

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch_Success(t *testing.T) {
	want := &State{
		Enabled:  true,
		Priority: "HIGH",
		Title:    "Production Incident",
		Message:  "Database connection failures detected",
		Sound:    "critical",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	f := NewFetcher(srv.URL)
	got, err := f.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Enabled != want.Enabled {
		t.Errorf("Enabled: got %v, want %v", got.Enabled, want.Enabled)
	}
	if got.Priority != want.Priority {
		t.Errorf("Priority: got %q, want %q", got.Priority, want.Priority)
	}
	if got.Title != want.Title {
		t.Errorf("Title: got %q, want %q", got.Title, want.Title)
	}
	if got.Message != want.Message {
		t.Errorf("Message: got %q, want %q", got.Message, want.Message)
	}
	if got.Sound != want.Sound {
		t.Errorf("Sound: got %q, want %q", got.Sound, want.Sound)
	}
}

func TestFetch_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := NewFetcher(srv.URL)
	_, err := f.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 response, got nil")
	}
}

func TestFetch_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	f := NewFetcher(srv.URL)
	_, err := f.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestFetch_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	f := NewFetcher(srv.URL)
	_, err := f.Fetch(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}
