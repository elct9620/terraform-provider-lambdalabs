package lambdalabs

type Client struct {
	ApiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey}
}
