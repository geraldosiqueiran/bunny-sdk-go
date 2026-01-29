package bunny

// PaginatedResponse represents a paginated API response.
type PaginatedResponse[T any] struct {
	Items       []T `json:"Items"`
	CurrentPage int `json:"CurrentPage"`
	TotalItems  int `json:"TotalItems"`
	ItemsPerPage int `json:"ItemsPerPage"`
}

// HasMore returns true if there are more items to fetch.
func (p *PaginatedResponse[T]) HasMore() bool {
	return p.CurrentPage*p.ItemsPerPage < p.TotalItems
}

// ListOptions specifies optional parameters for list operations.
type ListOptions struct {
	Page         int
	ItemsPerPage int
	Search       string
	OrderBy      string
}

// DefaultListOptions returns the default list options.
func DefaultListOptions() *ListOptions {
	return &ListOptions{
		Page:         1,
		ItemsPerPage: 100,
	}
}
