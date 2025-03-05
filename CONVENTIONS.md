# Coventions

## Testing

- Use standard library `testing` package for unit tests.
- Use `httptest` package for integration tests.

### Example

- Each function have a corresponding test function. e.g. `TestCreate` for `Create` function.
- Create `cases` variable to list all test cases.
- Expected output should be hard coded in the test case.

```go
func TestRetrieveInstance(t *testing.T) {
    cases := []struct {
        name string
        req  *lambdalabs.RetrieveInstanceRequest
        res  *lambdalabs.RetrieveInstanceResponse
        err  error
    }{
        {
            name: "success",
            req:   &lambdalabs.RetrieveInstanceRequest{
                Id: "123",
            },
            res:  &lambdalabs.Function{
                Data: &lambdalabs.Instance{
                    Id: "123",
                },
            },
            err:  nil,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"data": {"id": "123"}}`))
            }))
            defer mockServer.Close()

            client := lambdalabs.NewClient("mock-api-key", lambdalabs.WithBaseURL(mockServer.URL))
            res, err := client.RetrieveInstance(context.TODO(), c.req)

            if !reflect.DeepEqual(res, c.res) {
                t.Errorf("expected %v, got %v", c.res, res)
            }

            if err.Error() != c.err.Error() {
                t.Errorf("expected %v, got %v", c.err, err)
            }
        })
    }
}
```
