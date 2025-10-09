package http

import (
	"context"
	"net/http"

	ethproofs "github.com/nmvalera/go-utils/ethproofs/client"
)

func (c *Client) ListAWSPricing(ctx context.Context) ([]ethproofs.AWSInstance, error) {
	var resp []ethproofs.AWSInstance
	if err := c.do(ctx, http.MethodGet, "/aws-pricing-list", nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
