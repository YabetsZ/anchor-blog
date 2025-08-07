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

func (r *GeminiRepo) Generate(ctx context.Context, req entities.ContentRequest) (*entities.ContentResponse, error) {
	prompt := r.buildPrompt(req)
	if prompt == "" {
		return nil, errors.New("empty prompt generated")
	}

	payload := r.createPayload(prompt)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("payload marshaling failed: %w", err)
	}

	respBody, err := r.executeAPIRequest(ctx, body)
	if err != nil {
		return nil, err
	}

	return r.parseResponse(respBody)
}

func (r *GeminiRepo) buildPrompt(req entities.ContentRequest) string {
	return fmt.Sprintf(`Generate blog content about "%s" with:
	- Tone: %s
	- Word count: %d
	- Audience: %s
	- Safety: Avoid illegal content, explicit material, and harmful topics
	Output valid JSON without markdown using this structure:
	{
	"title": "Creative title",
	"outline": [
		{"heading": "H2 Heading", "key_points": ["bullet1", "bullet2"]}
	],
	"enhancements": {
		"seo_keywords": ["list", "of", "keywords"],
		"content_gaps": ["Missing elements"],
		"audience_tips": ["Audience-specific suggestions"]
	}
	}`, req.Topic, req.Tone, req.WordCount, req.Audience)
}

func (r *GeminiRepo) createPayload(prompt string) map[string]interface{} {
	safetySettings := make([]map[string]string, len(r.safetyConf.Categories))
	for i, category := range r.safetyConf.Categories {
		safetySettings[i] = map[string]string{
			"category":  category,
			"threshold": r.safetyConf.Threshold,
		}
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
			"stopSequences":   []string{"##"},
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

func (r *GeminiRepo) parseResponse(respBody []byte) (*entities.ContentResponse, error) {
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

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("API response parsing failed: %w", err)
	}

	if len(apiResp.Candidates) == 0 {
		return nil, AppError.ErrContentBlocked
	}

	candidate := apiResp.Candidates[0]

	if blocked, reasons := r.checkSafetyViolations(candidate.SafetyRatings); blocked {
		return nil, fmt.Errorf("%w: %v", AppError.ErrContentBlocked, reasons)
	}

	var content entities.ContentResponse
	if err := json.Unmarshal([]byte(candidate.Content.Parts[0].Text), &content); err != nil {
		return nil, fmt.Errorf("content JSON parsing failed: %w", err)
	}

	content.SafetyReport = r.buildSafetyReport(candidate.SafetyRatings)
	return &content, nil
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

func (r *GeminiRepo) buildSafetyReport(ratings []struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}) entities.SafetyReport {
	report := entities.SafetyReport{Safe: true}
	for _, rating := range ratings {
		if rating.Probability == "HIGH" || rating.Probability == "MEDIUM" {
			report.Blocked = true
			report.BlockReasons = append(report.BlockReasons,
				fmt.Sprintf("%s:%s", rating.Category, rating.Probability))
			report.Safe = false
		}
	}
	return report
}

type SafetyConfig struct {
	Threshold  string
	Categories []string
}
