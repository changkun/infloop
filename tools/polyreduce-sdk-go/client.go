// Copyright © 2022 The poly.red Authors. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package polyreduce

// ClientVersion defines polyreduce client version
const ClientVersion = "v0.0.1"

// DefaultEndpoint is the default endpoint of polyreduce service.
const DefaultEndpoint = "https://polyreduce.com"

// Client represents the client to interact the polyreduce service.
type Client struct {
	endpoint string
}

// NewClient creates a polyreduce client using default polyreduce endpoint.
func NewClient() *Client {
	return NewClientWithEndpoint(DefaultEndpoint)
}

// NewClientWithEndpoint creates a polyreduce client with a specific endpoint
func NewClientWithEndpoint(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// SetEndpoint sets the endpoint of polyreduce client
func (c *Client) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}
