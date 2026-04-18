package health_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourusername/vaultpipe/internal/health"
)

func newTestChecker(t *testing.T, handler http.Handler) *health.Checker {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("new vault client: %v", err)
	}
	return health.NewChecker(client)
}

// healthHandler returns an http.HandlerFunc that responds with the given status
// code and JSON body, setting Content-Type to application/json.
func healthHandler(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}
}

func TestCheck_Healthy(t *testing.T) {
	checker := newTestChecker(t, healthHandler(http.StatusOK,
		`{"initialized":true,"sealed":false,"standby":false,"version":"1.15.0"}`))
	s := checker.Check(context.Background())
	if !s.Healthy {
		t.Fatalf("expected healthy, got error: %v", s.Error)
	}
	if s.Version != "1.15.0" {
		t.Errorf("expected version 1.15.0, got %q", s.Version)
	}
}

func TestCheck_Sealed(t *testing.T) {
	checker := newTestChecker(t, healthHandler(503,
		`{"initialized":true,"sealed":true,"standby":false,"version":"1.15.0"}`))
	s := checker.Check(context.Background())
	if s.Healthy {
		t.Fatal("expected unhealthy for sealed vault")
	}
	if !s.Sealed {
		t.Fatal("expected sealed=true")
	}
}

func TestCheck_ConnectionRefused(t *testing.T) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = "http://127.0.0.1:19999"
	client, _ := vaultapi.NewClient(cfg)
	checker := health.NewChecker(client)
	s := checker.Check(context.Background())
	if s.Healthy {
		t.Fatal("expected unhealthy for unreachable server")
	}
	if s.Error == nil || !strings.Contains(s.Error.Error(), "health check failed") {
		t.Errorf("unexpected error: %v", s.Error)
	}
}
