package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dinoagera/AIChat/internal/domain"
	"github.com/google/uuid"
)

type ParsedEmergency struct {
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Priority    string  `json:"priority"`
	IssueType   string  `json:"issue_type"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
}
type LLMClient struct {
	urlAPI      string
	credentials string
	authURL     string
	httpClient  *http.Client
	log         *slog.Logger
	model       string
	accessToken string
	tokenExpiry time.Time
}

func NewLLMClient() {
	//todo:
}

func (c *LLMClient) ParseEmergencyText(ctx context.Context, text string) (*ParsedEmergency, error) {
	authCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := c.getAccessToken(authCtx); err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}
	prompt := c.buildEmergencyPrompt(text)
	messages := []domain.Message{
		{Role: "system", Content: domain.SystemPrompt},
		{Role: "user", Content: prompt},
	}

	requestBody := domain.Request{
		Model:             c.model,
		Messages:          messages,
		N:                 1,
		Stream:            false,
		MaxTokens:         1024,
		RepetitionPenalty: 1.0,
		UpdateInterval:    0,
		Temperature:       0.1,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.log.Error("Failed to marshal emergency parse request", "error", err)
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.urlAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		c.log.Error("Failed to create emergency parse request", "error", err)
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("Failed to send emergency parse request", "error", err)
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		c.log.Error("GigaChat API returned error", "status", resp.Status, "body", string(body))
		return nil, fmt.Errorf("GigaChat API error: %d - %s", resp.StatusCode, string(body))
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("Failed to read GigaChat response", "error", err)
		return nil, fmt.Errorf("read response: %w", err)
	}
	var aiResp domain.Response
	if err := json.Unmarshal(responseBody, &aiResp); err != nil {
		c.log.Error("Failed to unmarshal GigaChat response", "error", err, "body", string(responseBody))
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if len(aiResp.Choices) == 0 {
		c.log.Error("GigaChat returned no choices", "response", string(responseBody))
		return nil, fmt.Errorf("no choices in response")
	}
	content := strings.TrimSpace(aiResp.Choices[0].Message.Content)
	c.log.Debug("GigaChat raw response", "content", content)
	jsonContent := extractJSON(content)
	var result ParsedEmergency
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		c.log.Error("Failed to parse emergency JSON from LLM", "error", err, "content", content)
		return nil, fmt.Errorf("parse emergency JSON: %w", err)
	}
	if err := validateParsedEmergency(&result); err != nil {
		c.log.Warn("Parsed emergency data validation warning", "error", err, "data", result)
	}
	c.log.Info("Emergency text parsed successfully",
		"address", result.Address,
		"priority", result.Priority,
		"confidence", result.Confidence)
	return &result, nil
}

func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return content
	}
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		return content[start : end+1]
	}
	if strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]") {
		start := strings.Index(content, "[")
		end := strings.LastIndex(content, "]")
		if start != -1 && end != -1 && end > start {
			inner := content[start+1 : end]
			objStart := strings.Index(inner, "{")
			objEnd := strings.LastIndex(inner, "}")
			if objStart != -1 && objEnd != -1 && objEnd > objStart {
				return inner[objStart : objEnd+1]
			}
		}
	}
	return content
}

func validateParsedEmergency(p *ParsedEmergency) error {
	if strings.TrimSpace(p.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if p.Priority == "" {
		p.Priority = "normal"
	} else {
		validPriorities := map[string]bool{"critical": true, "high": true, "normal": true}
		if !validPriorities[p.Priority] {
			return fmt.Errorf("invalid priority '%s': must be critical/high/normal", p.Priority)
		}
	}
	if p.IssueType == "" {
		p.IssueType = "other"
	} else {
		validTypes := map[string]bool{
			"power_outage":      true,
			"equipment_failure": true,
			"fire":              true,
			"other":             true,
		}
		if !validTypes[p.IssueType] {
			return fmt.Errorf("invalid issue_type '%s'", p.IssueType)
		}
	}
	if p.Confidence == 0 {
		p.Confidence = 0.5
	} else if p.Confidence < 0 || p.Confidence > 1 {
		return fmt.Errorf("confidence out of range [0,1]: %f", p.Confidence)
	}
	if len(p.Description) > 200 {
		p.Description = p.Description[:200] + "..."
	}
	if (p.Latitude == 0) != (p.Longitude == 0) {
		p.Latitude = 0
		p.Longitude = 0
	}
	return nil
}
func (c *LLMClient) buildEmergencyPrompt(text string) string {
	return fmt.Sprintf(`Текст заявки от пользователя:
"""%s"""

Верни JSON:`, text)
}
func (l *LLMClient) getAccessToken(ctx context.Context) error {
	if !l.tokenExpiry.IsZero() && time.Now().Add(1*time.Minute).Before(l.tokenExpiry) {
		return nil
	}
	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_PERS")
	req, _ := http.NewRequestWithContext(ctx, "POST", l.authURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+l.credentials)
	req.Header.Set("RqUID", uuid.NewString())
	resp, _ := l.httpClient.Do(req)
	var authResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	json.NewDecoder(resp.Body).Decode(&authResp)
	l.accessToken = authResp.AccessToken
	l.tokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)
	return nil
}
