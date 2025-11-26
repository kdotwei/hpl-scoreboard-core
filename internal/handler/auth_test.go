package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kdotwei/hpl-scoreboard/internal/middleware"
	"github.com/stretchr/testify/assert"
)

// TestAuthMiddleware_Unauthorized verifies that requests without a token are rejected.
// This remains a "Red Phase" test because the middleware implementation is still missing.
func TestAuthMiddleware_Unauthorized(t *testing.T) {
	// 1. Setup: create a simple handler that mimics a protected API (e.g. uploading HPL files).
	// We do not have a router yet, so we mock a "fake" handler for now.
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Should return 200 once authentication passes
	})

	// Wrap the mockHandler with AuthMiddleware
	handlerToTest := middleware.AuthMiddleware(mockHandler)

	// 2. Execution: simulate sending a request without an Authorization header
	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload", nil)
	rr := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rr, req)

	// 3. Assertion: expect the system to return 401 Unauthorized.
	// Without the middleware it still returns 200, so the test stays Red.
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Should return 401 when no token provided")
}
