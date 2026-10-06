package api

import (
	"api-worker/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiClient struct {
	client  *http.Client
	baseURL string
}

func NewApiClient(baseURL string) *ApiClient {
	return &ApiClient{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (ac *ApiClient) Process(ctx context.Context, job model.Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ac.baseURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("external API returned status %d", resp.StatusCode)
	}
	return nil
}
