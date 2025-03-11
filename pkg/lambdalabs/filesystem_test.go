package lambdalabs_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
)

func TestListFileSystems(t *testing.T) {
	cases := []struct {
		name     string
		handler  http.HandlerFunc
		expected *lambdalabs.ListFileSystemsResponse
		err      error
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": []map[string]interface{}{
						{
							"id":          "398578a2336b49079e74043f0bd2cfe8",
							"name":        "my-filesystem",
							"mount_point": "/home/ubuntu/my-filesystem",
							"created":     "1970-01-01T00:00:00.000Z",
							"created_by": map[string]interface{}{
								"id":     "3da5a70a57a7422ea8a7203f98b2198b",
								"email":  "me@example.com",
								"status": "active",
							},
							"is_in_use": false,
							"region": map[string]interface{}{
								"name":        "us-west-1",
								"description": "California, USA",
							},
							"bytes_used": 0,
						},
					},
				})
			},
			expected: &lambdalabs.ListFileSystemsResponse{
				Data: []lambdalabs.FileSystem{
					{
						ID:         "398578a2336b49079e74043f0bd2cfe8",
						Name:       "my-filesystem",
						MountPoint: "/home/ubuntu/my-filesystem",
						Created:    "1970-01-01T00:00:00.000Z",
						CreatedBy: lambdalabs.User{
							ID:     "3da5a70a57a7422ea8a7203f98b2198b",
							Email:  "me@example.com",
							Status: "active",
						},
						IsInUse: false,
						Region: lambdalabs.Region{
							Name:        "us-west-1",
							Description: "California, USA",
						},
						BytesUsed: 0,
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
				if r.URL.Path != "/file-systems" {
					t.Errorf("Expected path %q, got %q", "/file-systems", r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method %q, got %q", http.MethodGet, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.ListFileSystems(context.Background())

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}

func TestDeleteFileSystem(t *testing.T) {
	cases := []struct {
		name     string
		request  *lambdalabs.DeleteFileSystemRequest
		handler  http.HandlerFunc
		expected *lambdalabs.DeleteFileSystemResponse
		err      error
	}{
		{
			name: "success",
			request: &lambdalabs.DeleteFileSystemRequest{
				ID: "398578a2336b49079e74043f0bd2cfe8",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": map[string]interface{}{
						"deleted_ids": []string{"398578a2336b49079e74043f0bd2cfe8"},
					},
				})
			},
			expected: &lambdalabs.DeleteFileSystemResponse{
				Data: struct {
					DeletedIDs []string `json:"deleted_ids"`
				}{
					DeletedIDs: []string{"398578a2336b49079e74043f0bd2cfe8"},
				},
			},
			err: nil,
		},
		{
			name: "bad request - filesystem in use",
			request: &lambdalabs.DeleteFileSystemRequest{
				ID: "398578a2336b49079e74043f0bd2cfe8",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"code":       "filesystems/filesystem-in-use",
						"message":    "The filesystem is currently in use by an instance",
						"suggestion": "Terminate the instance before deleting the filesystem",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "The filesystem is currently in use by an instance"},
		},
		{
			name: "unauthorized",
			request: &lambdalabs.DeleteFileSystemRequest{
				ID: "398578a2336b49079e74043f0bd2cfe8",
			},
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
			request: &lambdalabs.DeleteFileSystemRequest{
				ID: "398578a2336b49079e74043f0bd2cfe8",
			},
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
		{
			name: "not found",
			request: &lambdalabs.DeleteFileSystemRequest{
				ID: "nonexistent-id",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"code":       "global/object-does-not-exist",
						"message":    "Filesystem was not found.",
						"suggestion": "Check the filesystem ID and try again.",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "Filesystem was not found."},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/file-systems/" + c.request.ID
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %q, got %q", expectedPath, r.URL.Path)
				}
				if r.Method != http.MethodDelete {
					t.Errorf("Expected method %q, got %q", http.MethodDelete, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.DeleteFileSystem(context.Background(), c.request)

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}

func TestCreateFileSystem(t *testing.T) {
	cases := []struct {
		name     string
		request  *lambdalabs.CreateFileSystemRequest
		handler  http.HandlerFunc
		expected *lambdalabs.CreateFileSystemResponse
		err      error
	}{
		{
			name: "success",
			request: &lambdalabs.CreateFileSystemRequest{
				Name:   "my-filesystem",
				Region: "us-west-1",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Verify request body
				body, _ := io.ReadAll(r.Body)
				var req map[string]interface{}
				if err := json.Unmarshal(body, &req); err != nil {
					t.Errorf("Failed to unmarshal request body: %v", err)
					return
				}

				if req["name"] != "my-filesystem" || req["region"] != "us-west-1" {
					t.Errorf("Expected request with name 'my-filesystem' and region 'us-west-1', got %v", req)
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": map[string]interface{}{
						"id":          "398578a2336b49079e74043f0bd2cfe8",
						"name":        "my-filesystem",
						"mount_point": "/home/ubuntu/my-filesystem",
						"created":     "1970-01-01T00:00:00.000Z",
						"created_by": map[string]interface{}{
							"id":     "3da5a70a57a7422ea8a7203f98b2198b",
							"email":  "me@example.com",
							"status": "active",
						},
						"is_in_use": false,
						"region": map[string]interface{}{
							"name":        "us-west-1",
							"description": "California, USA",
						},
						"bytes_used": 0,
					},
				})
			},
			expected: &lambdalabs.CreateFileSystemResponse{
				Data: lambdalabs.FileSystem{
					ID:         "398578a2336b49079e74043f0bd2cfe8",
					Name:       "my-filesystem",
					MountPoint: "/home/ubuntu/my-filesystem",
					Created:    "1970-01-01T00:00:00.000Z",
					CreatedBy: lambdalabs.User{
						ID:     "3da5a70a57a7422ea8a7203f98b2198b",
						Email:  "me@example.com",
						Status: "active",
					},
					IsInUse: false,
					Region: lambdalabs.Region{
						Name:        "us-west-1",
						Description: "California, USA",
					},
					BytesUsed: 0,
				},
			},
			err: nil,
		},
		{
			name: "bad request - duplicate",
			request: &lambdalabs.CreateFileSystemRequest{
				Name:   "existing-filesystem",
				Region: "us-west-1",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"code":       "global/duplicate",
						"message":    "A file system with this name already exists",
						"suggestion": "Choose a different name for your file system",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "A file system with this name already exists"},
		},
		{
			name: "unauthorized",
			request: &lambdalabs.CreateFileSystemRequest{
				Name:   "my-filesystem",
				Region: "us-west-1",
			},
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
			request: &lambdalabs.CreateFileSystemRequest{
				Name:   "my-filesystem",
				Region: "us-west-1",
			},
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
				if r.URL.Path != "/file-systems" {
					t.Errorf("Expected path %q, got %q", "/file-systems", r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("Expected method %q, got %q", http.MethodPost, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.CreateFileSystem(context.Background(), c.request)

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}
