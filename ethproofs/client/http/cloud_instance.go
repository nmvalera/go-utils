package http

import (
	"context"
	"net/http"

	ethproofs "github.com/kkrt-labs/go-utils/ethproofs/client"
)

func (c *Client) ListCloudInstances(ctx context.Context) ([]ethproofs.CloudInstance, error) {
	var resp []ethproofs.CloudInstance
	if err := c.do(ctx, http.MethodGet, "/cloud-instances", nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
