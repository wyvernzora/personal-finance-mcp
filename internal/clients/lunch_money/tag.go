package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tags is a lookup map of tag ID to Tag.
// It is returned by ListTags for fast access by tag.Id.
type Tags map[int64]*Tag

// Tag represents a user-defined tag in the Lunch Money system.
// Each tag has an ID, a short name, an optional description, and an archive flag.
type Tag struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsArchived  bool   `json:"archived"`
}

// ListTags retrieves all tags from the Lunch Money API.
// It returns a Tags map keyed by tag ID or an error if the HTTP request or JSON unmarshalling fails.
func (c *client) ListTags(ctx context.Context) (Tags, error) {
	data, err := c.get(ctx, "/v1/tags", map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to call Lunch Money API: %w", err)
	}

	var response []*Tag
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	result := make(Tags)
	for _, tag := range response {
		result[tag.Id] = tag
	}

	return result, nil
}
