package lambdalabs

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransportAddsAuthorizationHeader(t *testing.T) {
	const expectedToken = "test-api-key"
	expectedAuthHeader := AuthorizationType + " " + expectedToken

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader != expectedAuthHeader {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuthHeader, authHeader)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(expectedToken, WithBaseUrl(server.URL))
	
	_, err := client.Get("/test", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
