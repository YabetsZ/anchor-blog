# Gemini Repository Documentation

## Overview

The `GeminiRepo` package provides an interface to Google's Gemini generative AI API for content generation. It handles prompt construction, API communication, and response processing to generate professional blog content in Markdown format.

## Features

- Structured content generation with professional formatting
- Safety filtering for content moderation
- Customizable generation parameters (tone, audience, length)
- Strict Markdown formatting enforcement
- Error handling for API failures and content violations

## Installation

```go
import "anchor-blog/internal/gemini"
```

## Initialization

Create a new Gemini repository instance:

```go
repo := gemini.NewGeminiRepo(apiKey, modelName)
```

Parameters:
- `apiKey`: Your Google Gemini API key
- `modelName`: The model to use (e.g., "gemini-pro")

## Methods

### `Generate(ctx context.Context, req entities.ContentRequest) (string, error)`

Generates blog content based on the provided request parameters.

**Parameters:**
- `ctx`: Context for request cancellation/timeout
- `req`: ContentRequest containing generation parameters

**Returns:**
- Generated Markdown content
- Error if generation fails

**Example:**
```go
content, err := repo.Generate(ctx, entities.ContentRequest{
    Topic:     "Remote Work Best Practices",
    Tone:      "Professional",
    Audience:  []string{"Managers", "Developers"},
    WordCount: 800,
    Scope:     "Focus on communication tools and productivity techniques",
})
```

## ContentRequest Structure

```go
type ContentRequest struct {
    Topic     string   // Primary subject of the content
    Tone      string   // Writing style (e.g., "Professional", "Casual")
    Audience  []string // Target reader demographics
    WordCount int      // Approximate length in words
    Scope     string   // Specific focus areas or constraints
}
```

## Safety Configuration

The repository includes built-in content safety checks that block:
- Hate speech
- Sexually explicit content
- Dangerous content
- Harassment

Content violating these thresholds will return an `ErrContentBlocked` error.

## Error Handling

Common errors:
- `ErrContentBlocked`: Generated content violated safety thresholds
- `ErrBadRequest`: Invalid input parameters
- `ErrInternalServer`: API communication failure

## Output Format

Generated content follows strict CommonMark Markdown rules:
- Proper header spacing
- Consistent list formatting
- Correct code block syntax
- Valid link formatting
- UNIX line endings

Example output structure:
```markdown
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
```

## Configuration

Default generation parameters:
- Temperature: 0.7 (creativity balance)
- Max tokens: 2000
- Top-P: 0.9
- Top-K: 40

These can be modified in the `createPayload` method if needed.

## Dependencies

- Google Gemini API
- Standard Go libraries: `net/http`, `encoding/json`, `context`