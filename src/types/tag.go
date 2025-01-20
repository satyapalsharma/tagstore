package types

import "time"

// Tag represents a tag that can be associated with files
type Tag struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// File represents a file tracked by TagStore
type File struct {
	ID        int64     `json:"id" db:"id"`
	Path      string    `json:"path" db:"path"`
	Hash      string    `json:"hash" db:"hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// FileTag represents the many-to-many relationship between files and tags
type FileTag struct {
	FileID    int64     `json:"file_id" db:"file_id"`
	TagID     int64     `json:"tag_id" db:"tag_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TagQuery represents a boolean query for searching files
// Supported operators: AND (space), OR (|), NOT (-)
// Example: "work -urgent | personal"
type TagQuery struct {
	Raw       string
	Positive  []string // tags that must be present
	Negative  []string // tags that must NOT be present
	Optional  []string // tags that are nice to have (OR groups)
}

// NewTagQuery parses a raw query string into a structured TagQuery
func NewTagQuery(raw string) *TagQuery {
	q := &TagQuery{Raw: raw}
	if raw == "" {
		return q
	}

	// Simple tokenizer: split by spaces, handle | for OR, - for NOT
	tokens := tokenize(raw)
	for _, token := range tokens {
		switch {
		case token == "|":
			// OR separator - next token goes to optional
			continue
		case strings.HasPrefix(token, "-"):
			q.Negative = append(q.Negative, strings.TrimPrefix(token, "-"))
		case token == "|":
			q.Optional = append(q.Optional, token)
		default:
			q.Positive = append(q.Positive, token)
		}
	}

	return q
}

func tokenize(s string) []string {
	var tokens []string
	var current strings.Builder
	for _, r := range s {
		switch r {
		case ' ', '|':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			if r == '|' {
				tokens = append(tokens, "|")
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}
