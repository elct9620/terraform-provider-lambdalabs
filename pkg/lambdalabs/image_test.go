package lambdalabs_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
)

func TestListImages(t *testing.T) {
	cases := []struct {
		name     string
		handler  http.HandlerFunc
		expected *lambdalabs.ListImagesResponse
		err      error
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": []map[string]interface{}{
						{
							"id":            "43336648-096d-4cba-9aa2-f9bb7727639d",
							"created_time":  "1970-01-01T00:00:00.000Z",
							"updated_time":  "1970-01-01T00:00:00.000Z",
							"name":          "ubuntu-24.04.01",
							"description":   "Ubuntu LTS",
							"family":        "ubuntu-lts",
							"version":       "24.04.01",
							"architecture":  "x86_64",
							"region": map[string]interface{}{
								"name":        "us-west-1",
								"description": "California, USA",
							},
						},
					},
				})
			},
			expected: &lambdalabs.ListImagesResponse{
				Data: []lambdalabs.Image{
					{
						ID:           "43336648-096d-4cba-9aa2-f9bb7727639d",
						CreatedTime:  "1970-01-01T00:00:00.000Z",
						UpdatedTime:  "1970-01-01T00:00:00.000Z",
						Name:         "ubuntu-24.04.01",
						Description:  "Ubuntu LTS",
						Family:       "ubuntu-lts",
						Version:      "24.04.01",
						Architecture: "x86_64",
						Region: lambdalabs.Region{
							Name:        "us-west-1",
							Description: "California, USA",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "unauthorized",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"code":       "global/invalid-api-key",
						"message":    "API key was invalid, expired, or deleted.",
						"suggestion": "Check your API key or create a new one, then try again.",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "API key was invalid, expired, or deleted."},
		},
		{
			name: "forbidden",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"code":       "global/account-inactive",
						"message":    "Your account is inactive.",
						"suggestion": "Make sure you have verified your email address and have a valid payment method. Contact Support if problems continue.",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "Your account is inactive."},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/images" {
					t.Errorf("Expected path %q, got %q", "/images", r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method %q, got %q", http.MethodGet, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.ListImages(context.Background())

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}
