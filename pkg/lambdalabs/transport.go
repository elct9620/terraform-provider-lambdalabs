package lambdalabs

import "net/http"

const (
	AuthorizationHeader = "Authorization"
	AuthorizationType   = "Bearer"
)

type Transport struct {
	apiKey string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(AuthorizationHeader, AuthorizationType+" "+t.apiKey)

	return http.DefaultTransport.RoundTrip(req)
}
