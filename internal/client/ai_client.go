package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CheckResult struct {
	Available       bool   `json:"available"`
	StatusCode      int    `json:"status_code"`
	LatencyMs       int64  `json:"latency_ms"`
	ErrorMessage    string `json:"error_message,omitempty"`
	ResponsePreview string `json:"response_preview,omitempty"`
	CheckTime       string `json:"check_time,omitempty"`
}

type responsesAPIRequest struct {
	Model  string `json:"model"`
	Input  string `json:"input"`
	Stream bool   `json:"stream"`
}

type responsesAPIResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`
	Status string `json:"status"`
	Output []struct {
		Type    string `json:"type"`
		ID      string `json:"id"`
		Role    string `json:"role"`
		Status  string `json:"status"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func CheckAPI(apiURL, token, model string) CheckResult {
	start := time.Now()
	result := CheckResult{
		StatusCode: 0,
		CheckTime:  time.Now().UTC().Format(time.RFC3339),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := strings.TrimRight(apiURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	endpoint := baseURL + "/responses"

	reqBody := responsesAPIRequest{
		Model:  model,
		Input:  "hi",
		Stream: false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		result.LatencyMs = time.Since(start).Milliseconds()
		result.ErrorMessage = "marshal request body failed: " + err.Error()
		return result
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		result.LatencyMs = time.Since(start).Milliseconds()
		result.ErrorMessage = "create request failed: " + err.Error()
		return result
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	result.LatencyMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Available = false
		result.ErrorMessage = "request failed: " + err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Available = false
		result.ErrorMessage = "read response body failed: " + err.Error()
		return result
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr errorEnvelope
		if err := json.Unmarshal(raw, &apiErr); err == nil && apiErr.Error.Message != "" {
			result.Available = false
			if apiErr.Error.Code != "" {
				result.ErrorMessage = fmt.Sprintf("%s: %s", apiErr.Error.Code, apiErr.Error.Message)
			} else {
				result.ErrorMessage = apiErr.Error.Message
			}
			return result
		}

		result.Available = false
		result.ErrorMessage = truncate(string(raw), 500)
		return result
	}

	var parsed responsesAPIResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		result.Available = false
		result.ErrorMessage = fmt.Sprintf(
			"invalid JSON body (content-type=%s): %v; raw=%s",
			resp.Header.Get("Content-Type"),
			err,
			truncate(string(raw), 500),
		)
		return result
	}

	text := extractOutputText(parsed)
	if text == "" {
		text = truncate(string(raw), 200)
	}

	result.Available = true
	result.ResponsePreview = truncate(text, 200)
	return result
}

func extractOutputText(resp responsesAPIResponse) string {
	var parts []string

	for _, out := range resp.Output {
		for _, c := range out.Content {
			if c.Type == "output_text" && strings.TrimSpace(c.Text) != "" {
				parts = append(parts, c.Text)
			}
		}
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
