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
	"log"
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
	log.Println(string(respBody))
	return r.parseResponse(respBody)
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
	audience := ""
	for _, word := range req.Audience {
		audience += " " + word
	}
	wordCount := req.WordCount
	if wordCount <= 0 {
		wordCount = 30
	}
	scope := sanitize(req.Scope, 1000)

	return fmt.Sprintf(`You are an expert blog writer and SEO specialist. Generate a high-quality blog post and return VALID JSON only (no markdown, no extra text, no explanation). The JSON must match this schema exactly:

{
  "title": "string",                     // 6-12 words, SEO-friendly
  "meta_description": "string",          // 1 sentence (max 155 characters)
  "outline": [
    {
      "heading": "H2 heading",
      "key_points": ["short bullet point", "..."]  // 1-4 items
    }
  ],
  "body": [
    {
      "heading": "H2 heading or intro",
      "paragraphs": ["paragraph1", "paragraph2"],
      "h3s": [
         {
           "subheading": "H3 heading",
           "bullets": ["bullet1", "bullet2"]
         }
      ]
    }
  ],
  "enhancements": {
    "seo_keywords": ["keyword1","keyword2"],
    "content_gaps": ["missing topic or angle to add"],
    "audience_tips": ["tips tailored for the audience"]
  }
}

User inputs:
- Topic: "%s"
- Tone: "%s"
- Audience: "%s"
- WordCount target: %d
- Scope/details: "%s"

Requirements:
1. Follow the schema above exactly. Return ONLY a single JSON object that validates to the schema.
2. Title: 6-12 words and include the main topic token(s).
3. Meta description: one concise sentence (<=155 chars).
4. Outline: at least 3 H2 sections; each H2 should include 1-4 key points.
5. Body: write paragraphs under each heading. Total body length should be approximately %d words (+/- 15%%). Use 2-4 short paragraphs per H2. Include at least 3 practical, actionable steps where relevant.
6. Use plain language appropriate for the audience. Provide localized examples if the scope mentions a location.
7. SEO: include a short list of 6-12 relevant SEO keywords in .
8. Safety: Avoid illegal content, explicit sexual content, instructions to harm, medical diagnosis or legal advice. Do not provide phone numbers, email addresses, or personal data.
9. Formatting: No markdown, no backticks, no explanation text — JSON ONLY.

If you cannot fulfill the request due to policy reasons, return JSON:
{ "error": "reason" }

Now generate the JSON-only response.`, topic, tone, audience, wordCount, scope, wordCount)
}

func (r *GeminiRepo) createPayload(prompt string) map[string]interface{} {
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
		return nil, fmt.Errorf("%w: no candidates returned", AppError.ErrContentBlocked)
	}

	candidate := apiResp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("empty content parts in API response")
	}

	if blocked, reasons := r.checkSafetyViolations(candidate.SafetyRatings); blocked {
		return nil, fmt.Errorf("%w: %v", AppError.ErrContentBlocked, strings.Join(reasons, ", "))
	}

	contentText := candidate.Content.Parts[0].Text
	cleaned := cleanMarkdownJSON(contentText)

	var temp struct {
		Title           string `json:"title"`
		MetaDescription string `json:"meta_description"`
		Outline         []struct {
			Heading   string   `json:"heading"`
			KeyPoints []string `json:"key_points"`
		} `json:"outline"`
		Body []struct {
			Heading    string   `json:"heading"`
			Paragraphs []string `json:"paragraphs"`
			H3s        []struct {
				Subheading string   `json:"subheading"`
				Bullets    []string `json:"bullets"`
			} `json:"h3s"`
		} `json:"body"`
		Enhancements struct {
			SEOKeywords  []string `json:"seo_keywords"`
			ContentGaps  []string `json:"content_gaps"`
			AudienceTips []string `json:"audience_tips"`
		} `json:"enhancements"`
		Error string `json:"error,omitempty"`
	}

	if err := json.Unmarshal([]byte(cleaned), &temp); err != nil {
		return nil, fmt.Errorf("content JSON parsing failed: %w\nOriginal: %s\nCleaned: %s",
			err, contentText, cleaned)
	}

	if temp.Error != "" {
		return nil, fmt.Errorf("%w: %s", AppError.ErrContentBlocked, temp.Error)
	}

	if len(temp.Title) == 0 || len(temp.MetaDescription) == 0 || len(temp.Outline) < 3 || len(temp.Body) < 3 {
		return nil, fmt.Errorf("invalid content structure: missing required fields")
	}

	content := entities.ContentResponse{
		Title:           temp.Title,
		MetaDescription: temp.MetaDescription,
		Outline:         make([]entities.Section, len(temp.Outline)),
		Body:            make([]entities.BodySection, len(temp.Body)),
		Enhancements: entities.Enhancements{
			SEOKeywords:  temp.Enhancements.SEOKeywords,
			ContentGaps:  temp.Enhancements.ContentGaps,
			AudienceTips: temp.Enhancements.AudienceTips,
		},
	}

	for i, o := range temp.Outline {
		content.Outline[i] = entities.Section{
			Heading:   o.Heading,
			KeyPoints: o.KeyPoints,
		}
	}

	for i, b := range temp.Body {
		content.Body[i] = entities.BodySection{
			Heading:     b.Heading,
			Paragraphs:  b.Paragraphs,
			Subsections: make([]entities.Subsection, len(b.H3s)),
		}
		for j, h3 := range b.H3s {
			content.Body[i].Subsections[j] = entities.Subsection{
				Subheading: h3.Subheading,
				Bullets:    h3.Bullets,
			}
		}
	}

	content.SafetyReport = r.buildSafetyReport(candidate.SafetyRatings)
	return &content, nil
}

func cleanMarkdownJSON(input string) string {
	input = strings.TrimPrefix(input, "```json")
	input = strings.TrimPrefix(input, "```")
	input = strings.TrimSuffix(input, "```")

	input = strings.ReplaceAll(input, "\\\"", "\"")
	input = strings.ReplaceAll(input, "\\n", "\n")
	input = strings.TrimSpace(input)

	if start := strings.Index(input, "{"); start > 0 {
		input = input[start:]
	}
	if end := strings.LastIndex(input, "}"); end >= 0 && end < len(input)-1 {
		input = input[:end+1]
	}

	return input
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
