package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OpenAPIEndpoint = "https://api.openai.com/v1/chat/completions"
	ModelName       = "gpt-4o-mini"
)

type OpenAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type OpenAPI struct {
	httpClient *http.Client
	cntxt      context.Context
	apiKey     string
}

// Constructor function
func NewOpenAPI(cntxt context.Context, apiKey string, httpClient *http.Client) *OpenAPI {
	o := &OpenAPI{
		cntxt:  cntxt,
		apiKey: apiKey,
	}

	if httpClient == nil {
		o.httpClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	} else {
		o.httpClient = httpClient
	}

	return o
}

// RESPONSE FUNCTION
func (o *OpenAPI) Query(systemPrompt string, prompt string) (OpenAPIResponse, error) {
	var response OpenAPIResponse
	if systemPrompt == " " {
		systemPrompt = "System Prompt is empty"
	}
	// byte slice, error
	bslice, err := json.Marshal(map[string]interface{}{
		"model": ModelName,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	})
	if err != nil {
		return response, err
	}
	// *http request, error
	req, err := http.NewRequestWithContext(o.cntxt, "POST", OpenAPIEndpoint, bytes.NewBuffer(bslice))
	if err != nil {
		return response, fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	// *http response, error
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()
	// byte slice,error
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("Error reading response: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("Error Unmarshaling response: %w", err)
	}

	if response.Error != nil {
		return response, fmt.Errorf("API Error: %s", response.Error.Message)
	}
	if len(response.Choices) == 0 {
		return response, fmt.Errorf("No choices returned from the API")
	}
	return response, nil
}
