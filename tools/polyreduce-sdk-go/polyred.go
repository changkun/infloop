// Copyright © 2022 The poly.red Authors. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package polyreduce

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type PolyredUploadInput struct {
	// ModelPath refers to an FBX file.
	ModelPath string
}
type PolyredUploadOutput struct {
	// The stored model ID that can be reused anytime in subsequent requests.
	ModelId string `json:"id,omitempty"`
	Message string `json:"msg,omitempty"`
}

// Upload uploads an given FBX model to the polyreduce service using
// plain polyred service.
func (c *Client) PolyredUpload(ctx context.Context, i *PolyredUploadInput) (*PolyredUploadOutput, error) {
	url := c.endpoint + "/api/v1/polyred/upload"

	if !strings.HasSuffix(strings.ToLower(i.ModelPath), ".fbx") {
		return nil, errors.New("only .FBX model is supported")
	}

	b, err := os.ReadFile(i.ModelPath)
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", i.ModelPath)
	if err != nil {
		return nil, err
	}
	part.Write(b)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.SetBasicAuth("way", "secret-pass")

	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &PolyredUploadOutput{}
	err = json.Unmarshal(data, output)
	return output, err
}

type PolyredConfigInput struct {
	ModelID        string
	ReductionRatio map[string]float64
}

type PolyredConfigOutput struct {
	// The stored model ID that can be reused anytime in subsequent requests.
	ModelId string `json:"id,omitempty"`
	Message string `json:"msg,omitempty"`
}

// PolyredConfig configs a simplification by providing the model ID and
// target reduction ratio. The ReductionRatio is a hash map that maps from
// object name to the target reduction ratio which allows multi-layer
// simplification.
//
// The configuration is allowed to call multiple times for a reconfiguration.
func (c *Client) PolyredConfig(ctx context.Context, i *PolyredConfigInput) error {
	url := c.endpoint + "/api/v1/polyred/config/" + i.ModelID
	b, err := json.Marshal(struct {
		Percent map[string]float64 `json:"percent"`
	}{
		Percent: i.ReductionRatio,
	})
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.SetBasicAuth("way", "secret-pass")
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	o := &PolyredConfigOutput{}
	err = json.Unmarshal(data, o)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send configuration: %s", o.Message)
	}
	return nil
}

type PolyredRunInput struct {
	ModelID string
}

type PolyredRunOutput struct {
	// The stored model ID that can be reused anytime in subsequent requests.
	ModelId string `json:"id,omitempty"`
	Message string `json:"msg,omitempty"`
}

// PolyredRun executes the simplification. The function blocks until the
// simplification is complete or server side error.
func (c *Client) PolyredRun(ctx context.Context, i *PolyredRunInput) error {
	url := c.endpoint + "/api/v1/polyred/run/" + i.ModelID
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	r.SetBasicAuth("way", "secret-pass")
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	o := &PolyredRunOutput{}
	err = json.Unmarshal(data, o)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to run: %s", o.Message)
	}
	return nil
}

type DownloadInput struct {
	ModelID string
	Path    string
}

// PolyredDownload downloads the result of a polygon reduction. The
// downloaded model is saved to the given path.
//
// The result may be different if the simplification was reconfigured and
// also being executed.
func (c *Client) PolyredDownload(ctx context.Context, i *DownloadInput) error {
	url := c.endpoint + "/api/v1/polyred/download/" + i.ModelID
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	r.SetBasicAuth("way", "secret-pass")
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(i.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
