package utils

import (
	"fmt"
	"strings"

	"tagstore/src/types"
)

// ParseTagQuery parses a raw query string into a structured query
func ParseTagQuery(raw string) (*types.TagQuery, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	query := types.NewTagQuery(raw)
	if query == nil {
		return nil, fmt.Errorf("failed to parse query")
	}

	return query, nil
}
