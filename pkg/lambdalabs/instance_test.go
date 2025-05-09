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

func TestRetrieveInstance(t *testing.T) {
	cases := []struct {
		name     string
		req      *lambdalabs.RetrieveInstanceRequest
		handler  http.HandlerFunc
		expected *lambdalabs.RetrieveInstanceResponse
		err      error
	}{
		{
			name: "success",
			req: &lambdalabs.RetrieveInstanceRequest{
				Id: "inst-123456",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": map[string]interface{}{
						"id":     "inst-123456",
						"name":   "test-instance",
						"ip":     "1.2.3.4",
						"status": "active",
					},
				})
			},
			expected: &lambdalabs.RetrieveInstanceResponse{
				Data: lambdalabs.Instance{
					ID:     "inst-123456",
					Name:   "test-instance",
					IP:     "1.2.3.4",
					Status: "active",
				},
			},
			err: nil,
		},
		{
			name: "not found",
			req: &lambdalabs.RetrieveInstanceRequest{
				Id: "inst-notexist",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"message": "Instance not found",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "Instance not found"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/instances/" + c.req.Id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %q, got %q", expectedPath, r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method %q, got %q", http.MethodGet, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.RetrieveInstance(context.Background(), c.req)

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}

func TestLaunchInstance(t *testing.T) {
	cases := []struct {
		name     string
		req      *lambdalabs.LaunchInstanceRequest
		handler  http.HandlerFunc
		expected *lambdalabs.LaunchInstanceResponse
		err      error
	}{
		{
			name: "success",
			req: &lambdalabs.LaunchInstanceRequest{
				RegionName:       "us-east-1",
				InstanceTypeName: "gpu-1x-a100",
				SSHKeyNames:      []string{"my-key"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": map[string]interface{}{
						"instance_ids": []string{"inst-123456"},
					},
				})
			},
			expected: &lambdalabs.LaunchInstanceResponse{
				Data: struct {
					IDs []string `json:"instance_ids"`
				}{
					IDs: []string{"inst-123456"},
				},
			},
			err: nil,
		},
		{
			name: "invalid instance type",
			req: &lambdalabs.LaunchInstanceRequest{
				RegionName:       "us-east-1",
				InstanceTypeName: "invalid-type",
				SSHKeyNames:      []string{"my-key"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"message": "Invalid instance type",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "Invalid instance type"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/instance-operations/launch" {
					t.Errorf("Expected path %q, got %q", "/instance-operations/launch", r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("Expected method %q, got %q", http.MethodPost, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.LaunchInstance(context.Background(), c.req)

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}

func TestListInstanceTypes(t *testing.T) {
	cases := []struct {
		name     string
		handler  http.HandlerFunc
		expected *lambdalabs.ListInstanceTypesResponse
		err      error
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": map[string]interface{}{
						"gpu_1x_gh200": map[string]interface{}{
							"instance_type": map[string]interface{}{
								"name":                "gpu_1x_gh200",
								"description":         "1x GH200 (96 GB)",
								"gpu_description":     "GH200 (96 GB)",
								"price_cents_per_hour": 149,
								"specs": map[string]interface{}{
									"vcpus":       64,
									"memory_gib":  432,
									"storage_gib": 4096,
									"gpus":        1,
								},
							},
							"regions_with_capacity_available": []map[string]interface{}{
								{
									"name":        "us-west-1",
									"description": "California, USA",
								},
							},
						},
					},
				})
			},
			expected: &lambdalabs.ListInstanceTypesResponse{
				Data: map[string]lambdalabs.InstanceTypeInfo{
					"gpu_1x_gh200": {
						InstanceType: lambdalabs.InstanceType{
							Name:               "gpu_1x_gh200",
							Description:        "1x GH200 (96 GB)",
							GPUDescription:     "GH200 (96 GB)",
							PriceCentsPerHour: 149,
							Specs: lambdalabs.InstanceTypeSpecs{
								VCPUs:      64,
								MemoryGiB:  432,
								StorageGiB: 4096,
								GPUs:       1,
							},
						},
						RegionsWithCapacityAvailable: []lambdalabs.Region{
							{
								Name:        "us-west-1",
								Description: "California, USA",
							},
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
						"message": "Unauthorized access",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "Unauthorized access"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/instance-types" {
					t.Errorf("Expected path %q, got %q", "/instance-types", r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method %q, got %q", http.MethodGet, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.ListInstanceTypes(context.Background())

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}

func TestTerminateInstance(t *testing.T) {
	cases := []struct {
		name     string
		req      *lambdalabs.TerminateInstanceRequest
		handler  http.HandlerFunc
		expected *lambdalabs.TerminateInstanceResponse
		err      error
	}{
		{
			name: "success",
			req: &lambdalabs.TerminateInstanceRequest{
				Ids: []string{"inst-123456"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"data": map[string]interface{}{
						"terminated_instances": []map[string]interface{}{
							{
								"id":     "inst-123456",
								"status": "terminated",
							},
						},
					},
				})
			},
			expected: &lambdalabs.TerminateInstanceResponse{
				Data: struct {
					TerminatedInstances []*lambdalabs.Instance `json:"terminated_instances"`
				}{
					TerminatedInstances: []*lambdalabs.Instance{
						{
							ID:     "inst-123456",
							Status: "terminated",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "instance not found",
			req: &lambdalabs.TerminateInstanceRequest{
				Ids: []string{"inst-notexist"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]interface{}{ // nolint:errcheck
					"error": map[string]string{
						"message": "Instance not found",
					},
				})
			},
			expected: nil,
			err:      &lambdalabs.Error{Message: "Instance not found"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/instance-operations/terminate" {
					t.Errorf("Expected path %q, got %q", "/instance-operations/terminate", r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("Expected method %q, got %q", http.MethodPost, r.Method)
				}

				c.handler(w, r)
			}))
			defer server.Close()

			client := lambdalabs.New("test-key", lambdalabs.WithBaseUrl(server.URL))
			result, err := client.TerminateInstance(context.Background(), c.req)

			if !reflect.DeepEqual(c.expected, result) {
				t.Errorf("Expected %+v, got %+v", c.expected, result)
			}

			if err != nil && c.err != nil && err.Error() != c.err.Error() {
				t.Errorf("Expected error %v, got %v", c.err, err)
			}
		})
	}
}
