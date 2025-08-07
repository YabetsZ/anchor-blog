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
		wordCount = 100
	}
	scope := sanitize(req.Scope, 1000)

	return fmt.Sprintf(`Generate a comprehensive blog post in MARKDOWN format following these requirements:

# Content Requirements
- Topic: %s
- Tone: %s
- Audience: %s
- Word Count: %d
- Scope: %s

# Formatting Rules
1. Use standard Markdown formatting
2. Include these sections:
   ## Title (H1)
   ### Meta Description (plain text)
   ## Outline (H2)
   - Bullet points of key sections
   ## Body Content (H2)
   - Detailed paragraphs
   ### Subsections (H3 as needed)
   - Actionable items as bullet points
   ## Enhancements (H2)
   - SEO Keywords
   - Content Gaps
   - Audience Tips

# Content Guidelines
- Provide 3-5 actionable tips
- Include 6-12 relevant SEO keywords
- Identify 2-3 content gaps
- Offer audience-specific advice
- Avoid any unsafe/prohibited content

# Example Structure
# Optimizing Developer Productivity

Meta description: Practical strategies to improve coding efficiency and focus for software engineers.

## Outline
- Time management techniques
- Tooling recommendations
- Team collaboration strategies

## Body Content
### Time Management
- Implement Pomodoro technique...
- Use time blocking...

### Recommended Tools
- VS Code extensions...
- CLI productivity tools...

## Enhancements
**SEO Keywords**: developer productivity, coding efficiency, time management  
**Content Gaps**: Comparison of IDEs, Remote pair programming tools  
**Audience Tips**: Adjust techniques for agile teams, Async communication practices

Now generate the content about: %s`, topic, tone, audience, wordCount, scope, topic)
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
