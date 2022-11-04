package filters

// Pagination is a pre-defined pagination struct, used within the listing endpoints.
type Pagination struct {
	Page         int `json:"page" yaml:"page"`
	ItemsPerPage int `json:"items_per_page" yaml:"items_per_page"`
	TotalCount   int `json:"total_count" yaml:"total_count"`
}

// GenerateOffset is used within limit queries, to generate proper pagination.
func GenerateOffset(pageNum int, itemPerPage int) int {
	return (pageNum - 1) * itemPerPage
}
