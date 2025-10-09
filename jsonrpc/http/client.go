package jsonrpchttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/nmvalera/go-utils/jsonrpc"
	comhttp "github.com/nmvalera/go-utils/net/http"
	comurl "github.com/nmvalera/go-utils/net/url"
)

// Client allows to connect to a JSON-RPC server
type Client struct {
	client autorest.Sender
}

// NewClient creates a new client capable of connecting to a JSON-RPC server
func NewClientFromClient(s autorest.Sender) *Client {
	c := &Client{
		client: s,
	}

	return c
}

// NewClient creates a new client capable of connecting to a JSON-RPC server
func NewClient(addr string, cfg *Config) (*Client, error) {
	u, err := comurl.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme for websocket connection: %s", u.Scheme)
	}

	httpc, err := comhttp.NewClient(cfg.HTTP)
	if err != nil {
		return nil, err
	}

	return NewClientFromClient(
		autorest.Client{
			Sender:           httpc,
			RequestInspector: comhttp.WithBaseURL(u),
		},
	), nil
}

func (c *Client) Call(ctx context.Context, r *jsonrpc.Request, res interface{}) error {
	return c.call(ctx, r, res)
}

// Call performs JSON-RPC call
func (c *Client) call(ctx context.Context, r *jsonrpc.Request, res interface{}) error {
	req, err := prepareCallRequest(ctx, r)
	if err != nil {
		msg, _ := json.Marshal(r)
		return autorest.NewErrorWithError(err, "jsonrpchttp.Client", fmt.Sprintf("Call(%v)", string(msg)), nil, "PrepareRequest")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		msg, _ := json.Marshal(r)
		return autorest.NewErrorWithError(err, "jsonrpchttp.Client", fmt.Sprintf("Call(%v)", string(msg)), resp, "Do")
	}

	err = inspectCallResponse(resp, res)
	if err != nil {
		msg, _ := json.Marshal(r)
		return autorest.NewErrorWithError(err, "jsonrpchttp.Client", fmt.Sprintf("Call(%v)", string(msg)), resp, "Inspect Response")
	}

	return nil
}

// ByUnmarshallingResponse marshall JSON-RPC request message into http.Request body
func prepareCallRequest(ctx context.Context, req *jsonrpc.Request) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithPath("/"),
		autorest.AsJSON(),
		autorest.WithJSON(req),
	).Prepare(newRequest(ctx))
}

func newRequest(ctx context.Context) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, "", "", http.NoBody)
	return req
}

func inspectCallResponse(resp *http.Response, res interface{}) error {
	msg := new(jsonrpc.ResponseMsg)
	err := autorest.Respond(
		resp,
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(msg),
		autorest.ByClosing(),
	)
	if err != nil {
		return err
	}

	return msg.Unmarshal(res)
}
