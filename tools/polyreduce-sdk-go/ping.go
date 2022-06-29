// Copyright © 2022 The poly.red Authors. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package polyreduce

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// PingInput is a a reserved structure
type PingInput struct {
}

// PingOutput is used for service health
type PingOutput struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	Message   string `json:"message"`
}

// Ping for polyreduce service health checking
func (c *Client) Ping(ctx context.Context) (*PingOutput, error) {
	url := c.endpoint + "/api/v1/ping"

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	o := &PingOutput{}
	err = json.Unmarshal(data, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
