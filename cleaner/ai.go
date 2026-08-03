package cleaner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func AIClean(records []Record, headers []string) ([]Record, error) {
	if len(records) == 0 {
		return records, nil
	}

	//Convert records to a readable table for Claude
	table := buildTable(records, headers)

	//prompt
	prompt := fmt.Sprintf(`You are a data cleaning assistant. 
Clean the following CSV data and return ONLY the cleaned CSV with the same headers.
Rules:
- Standardise inconsistent values (e.g "NYC", "New York City", "N.Y.C" should all become "New York City")
- Fix obvious typos
- Standardise date formats to YYYY-MM-DD
- Keep the same number of rows and columns
- Return ONLY the CSV data, no explanation

CSV Data:
%s`, table)

	// Call Claude
	cleaned, err := callClaude(prompt)
	if err != nil {
		return records, err
	}

	//  Parse Claude's response back into records
	reader := strings.NewReader(cleaned)
	cleanedRecords, _, err := ParseCSV(reader)
	if err != nil {
		return records, err
	}

	return cleanedRecords, nil
}

// buildTable converts records into a CSV string for the prompt
func buildTable(records []Record, headers []string) string {
	var buf bytes.Buffer

	buf.WriteString(strings.Join(headers, ",") + "\n")

	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = record[header]
		}
		buf.WriteString(strings.Join(row, ",") + "\n")
	}

	return buf.String()
}

// callClaude sends a prompt to Claude and returns the response
func callClaude(prompt string) (string, error) {
	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	authToken := os.Getenv("ANTHROPIC_AUTH_TOKEN")
	model := os.Getenv("ANTHROPIC_MODEL")

	if authToken == "" {
		return "", fmt.Errorf("ANTHROPIC_AUTH_TOKEN not set")
	}
	if baseURL == "" {
		return "", fmt.Errorf("ANTHROPIC_BASE_URL not set")
	}
	if model == "" {
		return "", fmt.Errorf("ANTHROPIC_MODEL not set")
	}

	// Build request body
	reqBody := AnthropicRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	// Marshal to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build the full URL
	url := baseURL + "v1/messages"

	// Create HTTP request
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", authToken)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return anthropicResp.Content[0].Text, nil
}
