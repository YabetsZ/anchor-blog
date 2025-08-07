package entities

import "github.com/go-playground/validator/v10"

// ContentRequest defines the input parameters for content generation
type ContentRequest struct {
	Topic     string   `json:"topic" validate:"required,min=3,max=100"`
	Tone      string   `json:"tone" validate:"oneof=professional casual persuasive"`
	WordCount uint     `json:"word_count" validate:"min=300,max=2000"`
	Audience  []string `json:"audience" validate:"required,dive,oneof=general developers marketers"`
	Scope     string   `json:"scope" validate:"max=500"`
}

// ContentResponse defines the complete generated content structure
type ContentResponse struct {
	Title           string        `json:"title" validate:"required,min=10,max=80"`
	MetaDescription string        `json:"meta_description" validate:"required,max=155"`
	Outline         []Section     `json:"outline" validate:"required,min=3"` // Minimum 3 sections
	Body            []BodySection `json:"body" validate:"required,min=1"`
	Enhancements    Enhancements  `json:"enhancements" validate:"required"`
	SafetyReport    SafetyReport  `json:"safety_report" validate:"required"`
	WordCount       uint          `json:"word_count" validate:"min=300"`
}

// Section represents a content outline item (H2 level)
type Section struct {
	Heading   string   `json:"heading" validate:"required,min=5,max=120"`
	KeyPoints []string `json:"key_points" validate:"required,min=1,max=4,dive,min=10,max=150"`
}

// BodySection contains the full content for each section
type BodySection struct {
	Heading     string       `json:"heading" validate:"required,min=5,max=120"`
	Paragraphs  []string     `json:"paragraphs" validate:"required,min=2,max=4,dive,min=50,max=300"`
	Subsections []Subsection `json:"h3s" validate:"dive"`
}

// Subsection represents H3 level content with actionable items
type Subsection struct {
	Subheading string   `json:"subheading" validate:"required,min=5,max=80"`
	Bullets    []string `json:"bullets" validate:"required,min=1,max=5,dive,min=10,max=120"`
}

// Enhancements contains SEO and content improvement data
type Enhancements struct {
	SEOKeywords  []string `json:"seo_keywords" validate:"required,min=5,max=12,dive,min=2,max=30"`
	ContentGaps  []string `json:"content_gaps" validate:"max=3,dive,min=10,max=100"`
	AudienceTips []string `json:"audience_tips" validate:"max=3,dive,min=10,max=120"`
}

// SafetyReport contains content moderation results
type SafetyReport struct {
	Blocked      bool     `json:"blocked"`
	BlockReasons []string `json:"block_reasons,omitempty" validate:"max=5"`
	Safe         bool     `json:"safe"`
}

// Validate performs struct validation
func (c *ContentResponse) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
