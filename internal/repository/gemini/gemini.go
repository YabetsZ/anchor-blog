package gemini

import (
	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GeminiRepo struct {
	apiKey     string
	baseURL    string
	model      string
	safetyConf SafetyConfig
}

func NewGeminiRepo(apiKey, model string) *GeminiRepo {
	return &GeminiRepo{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/",
		model:   model,
		safetyConf: SafetyConfig{
			Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			Categories: []string{
				"HARM_CATEGORY_DANGEROUS",
				"HARM_CATEGORY_HARASSMENT",
				"HARM_CATEGORY_HATE_SPEECH",
				"HARM_CATEGORY_SEXUALLY_EXPLICIT",
				"HARM_CATEGORY_DANGEROUS_CONTENT",
			},
		},
	}
}

func (r *GeminiRepo) Generate(ctx context.Context, req entities.ContentRequest) (string, error) {
	prompt := r.buildPrompt(req)
	if prompt == "" {
		return "", errors.New("empty prompt generated")
	}

	payload := r.createPayload(prompt)
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("payload marshaling failed: %w", err)
	}

	respBody, err := r.executeAPIRequest(ctx, body)
	if err != nil {
		return "", err
	}

	// Parse and extract the Markdown content
	markdown, err := r.extractMarkdownResponse(respBody)
	if err != nil {
		return "", fmt.Errorf("failed to extract markdown: %w", err)
	}

	return markdown, nil
}

func (r *GeminiRepo) extractMarkdownResponse(respBody []byte) (string, error) {
	// Define response structure to extract the text part
	type apiCandidate struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	}

	type apiResponse struct {
		Candidates []apiCandidate `json:"candidates"`
	}

	// Parse the outer API response
	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("API response parsing failed: %w", err)
	}

	if len(apiResp.Candidates) == 0 {
		return "", AppError.ErrContentBlocked
	}

	candidate := apiResp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return "", errors.New("empty content parts in API response")
	}

	// Check safety violations
	if blocked, reasons := r.checkSafetyViolations(candidate.SafetyRatings); blocked {
		return "", fmt.Errorf("%w: %v", AppError.ErrContentBlocked, strings.Join(reasons, ", "))
	}

	// Return the raw Markdown content
	return candidate.Content.Parts[0].Text, nil
}

func sanitize(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	if maxLen > 0 && len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}

func (r *GeminiRepo) buildPrompt(req entities.ContentRequest) string {
	topic := sanitize(req.Topic, 300)
	tone := sanitize(req.Tone, 50)
	if tone == "" {
		tone = "informative"
	}
	audience := strings.Join(req.Audience, ", ")
	wordCount := req.WordCount
	if wordCount <= 0 {
		wordCount = 1000
	}
	scope := sanitize(req.Scope, 1000)

	return fmt.Sprintf(`Generate a professional blog post in STRICT CommonMark Markdown format. Follow these rules EXACTLY:

# FORMATTING RULES
1. Use UNIX line endings (LF only, no CRLF)
2. Exactly one blank line between sections
3. No trailing whitespace on any line
4. Headers must have exactly one space after #
5. Lists must use hyphen with one space ("- item")
6. Code blocks use triple backticks with language
7. Links must use [text](url) format

# REQUIRED STRUCTURE
# [Title: 6-12 Words]

[1-sentence meta description under 155 chars]

## Outline
- [Main point 1]
- [Main point 2]
- [Main point 3]

## Introduction
[2-3 paragraphs introducing topic]

## [Section 1]
[2-4 paragraphs]

### [Subsection]
- [Actionable tip 1]
- [Actionable tip 2]

## Enhancements
**SEO Keywords**: [comma-separated terms]  
**Content Gaps**: [missing aspects]  
**Audience Tips**: [tailored suggestions]

# CONTENT PARAMETERS
Topic: %s
Tone: %s
Audience: %s
Length: %d words (±15%%)
Scope: %s

# EXAMPLE OUTPUT
# Effective Remote Team Management

Meta description: Proven strategies to maintain productivity and collaboration in distributed teams.

## Outline
- Communication protocols
- Productivity tools
- Team building activities

## Introduction
Managing remote teams requires...

## Communication Protocols
Establish clear expectations...

### Best Practices
- Use async video updates
- Document all decisions
- Set core overlap hours

## Enhancements
**SEO Keywords**: remote work, team management, async communication  
**Content Gaps**: Timezone management tools, Security considerations  
**Audience Tips**: Managers: Schedule regular 1:1s, Developers: Use focus timers

Now generate content about: %s`,
		topic, tone, audience, wordCount, scope, topic)
}

func (r *GeminiRepo) createPayload(prompt string) map[string]interface{} {
	// Use the exact enum values expected by Gemini API
	safetySettings := []map[string]string{
		{
			"category":  "HARM_CATEGORY_HATE_SPEECH",
			"threshold": "BLOCK_MEDIUM_AND_ABOVE",
		},
		{
			"category":  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
			"threshold": "BLOCK_MEDIUM_AND_ABOVE",
		},
		{
			"category":  "HARM_CATEGORY_DANGEROUS_CONTENT",
			"threshold": "BLOCK_MEDIUM_AND_ABOVE",
		},
		{
			"category":  "HARM_CATEGORY_HARASSMENT",
			"threshold": "BLOCK_MEDIUM_AND_ABOVE",
		},
	}

	return map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{"text": prompt},
				},
			},
		},
		"safetySettings": safetySettings,
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 2000,
			"temperature":     0.7,
			"topP":            0.9,
			"topK":            40,
		},
	}
}

func (r *GeminiRepo) executeAPIRequest(ctx context.Context, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s:generateContent?key=%s", r.baseURL, r.model, r.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(errorBody))
	}

	return io.ReadAll(resp.Body)
}

func (r *GeminiRepo) checkSafetyViolations(ratings []struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}) (bool, []string) {
	violations := []string{}
	for _, rating := range ratings {
		if rating.Probability == "HIGH" || rating.Probability == "MEDIUM" {
			violations = append(violations,
				fmt.Sprintf("%s:%s", rating.Category, rating.Probability))
		}
	}
	return len(violations) > 0, violations
}

type SafetyConfig struct {
	Threshold  string
	Categories []string
}
