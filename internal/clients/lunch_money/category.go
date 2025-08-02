package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
)

type Categories map[int64]*Category

func (c Categories) Get(id int64) *Category {
	if cat, ok := c[id]; ok {
		return cat
	}
	return nil
}

// Category represents a transaction category as returned by the Lunch Money API.
// It includes the category’s ID, name, description, ordering, income/budget flags, archive status,
// timestamps for creation/update/archive, and any nested child categories.
type Category struct {
	Id                int64       `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	Order             int64       `json:"order"`
	IsIncome          bool        `json:"is_income"`
	ExcludeFromBudget bool        `json:"exclude_from_budget"`
	ExcludeFromTotals bool        `json:"exclude_from_totals"`
	IsArchived        bool        `json:"is_archived"`
	ArchivedOn        string      `json:"archived_on"`
	UpdatedAt         string      `json:"updated_at"`
	CreatedAt         string      `json:"created_at"`
	Children          []*Category `json:"children,omitempty"`
}

type listCategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

// ListCategories fetches all categories from the Lunch Money API in nested format ("format=nested").
// It returns a slice of Category or an error if the HTTP request or JSON unmarshalling fails.
func (c *client) ListCategories(ctx context.Context) (Categories, error) {
	data, err := c.get(ctx, "/v1/categories", map[string]string{
		"format": "nested",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call Lunch Money API: %w", err)
	}

	var response listCategoriesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	result := make(map[int64]*Category, len(response.Categories))
	indexCategories(result, response.Categories)

	return result, nil
}

func indexCategories(m map[int64]*Category, cats []*Category) {
	for _, c := range cats {
		m[c.Id] = c
		indexCategories(m, c.Children)
	}
}
