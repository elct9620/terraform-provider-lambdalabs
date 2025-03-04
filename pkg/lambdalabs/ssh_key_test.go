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

func TestListSshKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &lambdalabs.ListSshKeysResponse{
			Data: []lambdalabs.SshKey{
				{
					Id:        "ddf9a910ceb744a0bb95242cbba6cb50",
					Name:      "my-public-key",
					PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICN+lJwsONkwrdsSnQsu1ydUkIuIg5oOC+Eslvmtt60T noname",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ssh-keys" {
				t.Errorf("Expected path %q, got %q", "/ssh-keys", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("Expected method %q, got %q", http.MethodGet, r.Method)
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expected); err != nil {
				t.Fatal(err)
			}
		}))
		defer server.Close()

		client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
		result, err := client.ListSshKeys(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
				"error": map[string]string{
					"code":       "global/invalid-api-key",
					"message":    "API key was invalid, expired, or deleted.",
					"suggestion": "Check your API key or create a new one, then try again.",
				},
			})
		}))
		defer server.Close()

		client := lambdalabs.New("invalid-key", lambdalabs.WithBaseUrl(server.URL))
		_, err := client.ListSshKeys(context.Background())
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCreateSshKey(t *testing.T) {
	t.Run("success with existing public key", func(t *testing.T) {
		publicKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICN+lJwsONkwrdsSnQsu1ydUkIuIg5oOC+Eslvmtt60T noname"
		expected := &lambdalabs.CreateSshKeyResponse{
			Data: lambdalabs.SshKey{
				Id:        "ddf9a910ceb744a0bb95242cbba6cb50",
				Name:      "my-public-key",
				PublicKey: publicKey,
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ssh-keys" {
				t.Errorf("Expected path %q, got %q", "/ssh-keys", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("Expected method %q, got %q", http.MethodPost, r.Method)
			}

			var req lambdalabs.CreateSshKeyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatal(err)
			}
			if req.Name != "my-public-key" {
				t.Errorf("Expected name %q, got %q", "my-public-key", req.Name)
			}
			if *req.PublicKey != publicKey {
				t.Errorf("Expected public key %q, got %q", publicKey, *req.PublicKey)
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expected); err != nil {
				t.Fatal(err)
			}
		}))
		defer server.Close()

		client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
		result, err := client.CreateSshKey(context.Background(), &lambdalabs.CreateSshKeyRequest{
			Name:      "my-public-key",
			PublicKey: &publicKey,
		})
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("success generate new key pair", func(t *testing.T) {
		expected := &lambdalabs.CreateSshKeyResponse{
			Data: lambdalabs.SshKey{
				Id:         "ddf9a910ceb744a0bb95242cbba6cb50",
				Name:       "new-key",
				PublicKey:  "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICN+lJwsONkwrdsSnQsu1ydUkIuIg5oOC+Eslvmtt60T noname",
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEoAIBAAKCAQEAmTi0yMd35HkIKXgEAVLb14fE094YL5qgGqS5ayq9SHi72mlf\n-----END RSA PRIVATE KEY-----\n",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ssh-keys" {
				t.Errorf("Expected path %q, got %q", "/ssh-keys", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("Expected method %q, got %q", http.MethodPost, r.Method)
			}

			var req lambdalabs.CreateSshKeyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatal(err)
			}
			if req.Name != "new-key" {
				t.Errorf("Expected name %q, got %q", "new-key", req.Name)
			}
			if req.PublicKey != nil {
				t.Error("Expected nil public key")
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expected); err != nil {
				t.Fatal(err)
			}
		}))
		defer server.Close()

		client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
		result, err := client.CreateSshKey(context.Background(), &lambdalabs.CreateSshKeyRequest{
			Name: "new-key",
		})
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("bad request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
				"error": map[string]string{
					"code":       "global/invalid-parameters",
					"message":    "Invalid request data.",
					"suggestion": "Check your request parameters",
				},
			})
		}))
		defer server.Close()

		client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
		_, err := client.CreateSshKey(context.Background(), &lambdalabs.CreateSshKeyRequest{})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestDeleteSshKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ssh-keys/test-id" {
				t.Errorf("Expected path %q, got %q", "/ssh-keys/test-id", r.URL.Path)
			}
			if r.Method != http.MethodDelete {
				t.Errorf("Expected method %q, got %q", http.MethodDelete, r.Method)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
				"data": map[string]interface{}{},
			})
		}))
		defer server.Close()

		client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
		err := client.DeleteSshKey(context.Background(), &lambdalabs.DeleteSshKeyRequest{
			Id: "test-id",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
				"error": map[string]string{
					"code":    "global/not-found",
					"message": "SSH key not found",
				},
			})
		}))
		defer server.Close()

		client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
		err := client.DeleteSshKey(context.Background(), &lambdalabs.DeleteSshKeyRequest{
			Id: "non-existent-id",
		})
		if err == nil {
			t.Fatal("Expected error for non-existent SSH key")
		}
		if err.Error() != "SSH key not found" {
			t.Errorf("Expected error %v, got %v", "SSH key not found", err)
		}
	})
}
