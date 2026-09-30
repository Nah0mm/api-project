package api

import (
	"api-worker/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIClient struct {
	client  *http.Client
	baseURL string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (ac *APIClient) Process(ctx context.Context, job model.Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("Error marshalling job\n%v\n", err)

	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ac.baseURL,
		bytes.NewReader(body),
	)
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
