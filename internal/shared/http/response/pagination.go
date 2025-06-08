package response

type meta struct {
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	PerPage     int  `json:"per_page"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
}

type pagination[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
	Meta    meta   `json:"meta"`
}

func NewPagination[T any](status, message string, data []T, meta meta) pagination[T] {
	if len(data) == 0 {
		data = []T{}
	}

	return pagination[T]{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func NewMeta(currentPage, perPage, totalItems int) meta {
	if totalItems <= 0 {
		return meta{
			CurrentPage: 1,
			TotalPages:  1,
			TotalItems:  0,
			PerPage:     perPage,
			NextPage:    nil,
			PrevPage:    nil,
		}
	}

	var (
		prevPage *int
		nextPage *int
	)

	totalPages := (totalItems + perPage - 1) / perPage

	if totalPages > 1 && currentPage < totalPages {
		nextPageValue := currentPage + 1
		nextPage = &nextPageValue
	}

	if currentPage > 1 {
		prevPageValue := currentPage - 1
		prevPage = &prevPageValue
	}

	return meta{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		PerPage:     perPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}
