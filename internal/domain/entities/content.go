package entities

type ContentRequest struct {
	Topic     string   `json:"topic" validate:"required,min=3,max=100"`
	Tone      string   `json:"tone" validate:"oneof=professional casual persuasive"`
	WordCount uint     `json:"word_count" validate:"min=300,max=2000"`
	Audience  []string `json:"audience" validate:"oneof=general developers marketers"`
}

type ContentResponse struct {
	Title        string       `json:"title"`
	Outline      []Section    `json:"outline"`
	Enhancements Enhancements `json:"enhancements"`
	SafetyReport SafetyReport `json:"safety_report"` // Added safety metadata
}

type Section struct {
	Heading   string   `json:"heading"`
	KeyPoints []string `json:"key_points"`
}

type Enhancements struct {
	SEOKeywords  []string `json:"seo_keywords"`
	ContentGaps  []string `json:"content_gaps"`
	AudienceTips []string `json:"audience_tips"`
}

type SafetyReport struct {
	Blocked      bool     `json:"blocked"`
	BlockReasons []string `json:"block_reasons,omitempty"`
	Safe         bool     `json:"safe"`
}
