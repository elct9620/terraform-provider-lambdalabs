package lambdalabs

import (
	"context"
	"encoding/json"
)

// ListImagesResponse represents the response from the List Images API
type ListImagesResponse struct {
	Data []Image `json:"data"`
}

// ListImages retrieves all available images
func (c *Client) ListImages(ctx context.Context) (*ListImagesResponse, error) {
	resp, err := c.Get(ctx, "/images", nil)
	if err != nil {
		return nil, err
	}

	var res ListImagesResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
