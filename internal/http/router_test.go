package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiauth "github.com/sigdown/kartograf-backend-go/internal/auth"
)

func TestNewRouterRegistersHealthRoute(t *testing.T) {
	router := NewRouter(Services{
		Tokens: apiauth.NewTokenManager("secret", time.Hour),
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
