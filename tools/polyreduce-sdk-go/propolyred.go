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

type ProPolyredUploadInput struct {
	// ModelPath refers to an FBX file.
	ModelPath string
}
type ProPolyredUploadOutput struct {
	// The session ID that can be reused anytime in subsequent requests.
	// The session ID also represents the ID of the root (initial) model.
	SessionId string `json:"id,omitempty"`
	Message   string `json:"msg,omitempty"`
}

func (c *Client) ProPolyredUpload(ctx context.Context, i *ProPolyredUploadInput) (*ProPolyredUploadOutput, error) {
	url := c.endpoint + "/api/v1/propolyred/upload"

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
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &ProPolyredUploadOutput{}
	err = json.Unmarshal(data, output)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(output.Message)
	}
	return output, err
}

type ProPolyredRunInput struct {
	SessionId string
}
type ProPolyredRunOutput struct {
	Phases         []string `json:"ids,omitempty"`
	AssumedOptimal float64  `json:"optimal,omitempty"`
	Message        string   `json:"msg,omitempty"`
}

func (c *Client) ProPolyredRun(ctx context.Context, i *ProPolyredRunInput) (*ProPolyredRunOutput, error) {
	url := fmt.Sprintf("%s/api/v1/propolyred/run/%s", c.endpoint, i.SessionId)
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	r.SetBasicAuth("way", "secret-pass")

	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to send the request: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	output := &ProPolyredRunOutput{}
	err = json.Unmarshal(data, output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the run output: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(output.Message)
	}
	return output, nil
}

type ProPolyredDownloadInput struct {
	SessionId, PhaseId string
	Path               string
}

func (c *Client) ProPolyredDownload(ctx context.Context, i *ProPolyredDownloadInput) error {
	url := fmt.Sprintf("%s/api/v1/propolyred/download/%s/%s", c.endpoint, i.SessionId, i.PhaseId)
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

type ProPolyredInspectInput struct {
	SessionId string
}
type ProPolyredInspectOutput struct {
	Unevaluated []string `json:"ids,omitempty"`
	Message     string   `json:"msg,omitempty"`
}

// ProPolyredInspect returns a list of model IDs that are not yet evaluated.
func (c *Client) ProPolyredInspect(ctx context.Context, i *ProPolyredInspectInput) (*ProPolyredInspectOutput, error) {
	url := c.endpoint + "/api/v1/propolyred/evaluate/" + i.SessionId

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
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

	output := &ProPolyredInspectOutput{}
	err = json.Unmarshal(data, output)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(output.Message)
	}

	return output, err
}

type ProPolyredEvaluateInput struct {
	SessionId string
	Rating    map[string]float64
}

type ProPolyredEvaluateOutput struct {
	Message string `json:"msg,omitempty"`
}

func (c *Client) ProPolyredEvaluate(ctx context.Context, i *ProPolyredEvaluateInput) error {
	url := c.endpoint + "/api/v1/propolyred/evaluate/" + i.SessionId

	b, err := json.Marshal(i.Rating)
	if err != nil {
		return fmt.Errorf("failed to marshal rating: %w", err)
	}

	r, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	r.SetBasicAuth("way", "secret-pass")

	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	output := &ProPolyredEvaluateOutput{}
	err = json.Unmarshal(data, output)
	if err != nil {
		return fmt.Errorf("failed to parse evaluate output: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(output.Message)
	}
	return nil
}

type ProPolyredResetInput struct {
	SessionId string
}
type ProPolyredResetOutput struct {
	// The session ID that can be reused anytime in subsequent requests.
	// The session ID also represents the ID of the root (initial) model.
	SessionId string `json:"id,omitempty"`
	Message   string `json:"msg,omitempty"`
}

func (c *Client) ProPolyredReset(ctx context.Context, i *ProPolyredResetInput) (*ProPolyredResetOutput, error) {
	url := c.endpoint + "/api/v1/propolyred/reset/" + i.SessionId
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
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

	output := &ProPolyredResetOutput{}
	err = json.Unmarshal(data, output)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(output.Message)
	}
	return output, err
}

type ProPolyredCopyInput struct {
	SessionId string
}
type ProPolyredCopyOutput struct {
	// The session ID that can be reused anytime in subsequent requests.
	// The session ID also represents the ID of the root (initial) model.
	SessionId string `json:"id,omitempty"`
	Message   string `json:"msg,omitempty"`
}

func (c *Client) ProPolyredCopy(ctx context.Context, i *ProPolyredCopyInput) (*ProPolyredCopyOutput, error) {
	url := c.endpoint + "/api/v1/propolyred/copy/" + i.SessionId
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	r.SetBasicAuth("way", "secret-pass")

	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &ProPolyredCopyOutput{}
	err = json.Unmarshal(data, output)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(output.Message)
	}
	return output, err
}
