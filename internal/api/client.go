package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	leetcodeGraphQLEndpoint = "https://leetcode.com/graphql/"
	defaultTimeout          = 30 * time.Second
)

// Client is a LeetCode GraphQL API client
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// GraphQLRequest is the structure of a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// doRequest executes a GraphQL request and returns the response body
// Since no graphql schema is provided, using no graphql client library
func (c *Client) doRequest(query string, variables map[string]interface{}) ([]byte, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", leetcodeGraphQLEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

func (c *Client) SearchQuestions(keyword string, limit, skip int) (*ProblemsetResponse, error) {
	variables := map[string]interface{}{
		"searchKeyword": keyword,
		"limit":         limit,
		"skip":          skip,
	}

	body, err := c.doRequest(ProblemsetPanelQuestionListQuery, variables)
	if err != nil {
		return nil, err
	}

	var response ProblemsetResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetQuestionDetail(titleSlug string) (*QuestionDetailResponse, error) {
	variables := map[string]interface{}{
		"titleSlug": titleSlug,
	}

	body, err := c.doRequest(QuestionDetailQuery, variables)
	if err != nil {
		return nil, err
	}

	var response QuestionDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
