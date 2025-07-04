package agents

import (
	"context"
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
