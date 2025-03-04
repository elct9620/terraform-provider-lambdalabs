package lambdalabs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
)

func TestTransportAddsAuthorizationHeader(t *testing.T) {
	const expectedToken = "test-api-key"
	expectedAuthHeader := lambdalabs.AuthorizationType + " " + expectedToken

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(lambdalabs.AuthorizationHeader)
		if authHeader != expectedAuthHeader {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuthHeader, authHeader)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := lambdalabs.New(expectedToken, lambdalabs.WithBaseUrl(server.URL))

	_, err := client.Get(context.TODO(), "/test", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
