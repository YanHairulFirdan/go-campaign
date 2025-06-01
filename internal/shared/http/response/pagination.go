package response

import "errors"

type DataRetrieverFunc func() ([]any, error)
type CountRetrieverFunc func() (int, error)

type paginationBuilder[T any] struct {
	PerPage        int
	CurrentPage    int
	DataRetriever  func() ([]T, error)
	EmptyData      []T
	CountRetriever CountRetrieverFunc
}

func NewPaginationBuilder[T any](
	perPage, currentPage int,
	dataRetriever func() ([]T, error),
	countRetriever CountRetrieverFunc,
) *paginationBuilder[T] {
	return &paginationBuilder[T]{
		PerPage:        perPage,
		CurrentPage:    currentPage,
		DataRetriever:  dataRetriever,
		CountRetriever: countRetriever,
		EmptyData:      []T{},
	}
}

func (pb *paginationBuilder[T]) Build() (*pagination[T], error) {
	if pb.DataRetriever == nil || pb.CountRetriever == nil {
		return nil, errors.New("data retriever or count retriever is not set")
	}

	data, err := pb.DataRetriever()
	if err != nil {
		return nil, err
	}

	count, err := pb.CountRetriever()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		data = pb.EmptyData
	}

	meta := pb.GetMeta(count)

	return &pagination[T]{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    data,
		Meta:    *meta,
	}, nil
}

func (pb *paginationBuilder[T]) GetMeta(totalCount int) *meta {
	if totalCount <= 0 {
		return &meta{
			CurrentPage: 1,
			TotalPages:  1,
			TotalItems:  0,
			PerPage:     pb.PerPage,
			NextPage:    nil,
			PrevPage:    nil,
		}
	}

	var (
		prevPage *int
		nextPage *int
	)

	totalPages := (totalCount + pb.PerPage - 1) / pb.PerPage

	if totalPages > 1 && pb.CurrentPage < totalPages {
		nextPageValue := pb.CurrentPage + 1
		nextPage = &nextPageValue
	}

	if pb.CurrentPage > 1 {
		prevPageValue := pb.CurrentPage - 1
		prevPage = &prevPageValue
	}

	return &meta{
		CurrentPage: pb.CurrentPage,
		TotalPages:  totalPages,
		TotalItems:  totalCount,
		PerPage:     pb.PerPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

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
	return pagination[T]{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func NewMeta(currentPage, totalPages, totalItems, perPage int, nextPage, prevPage *int) meta {
	return meta{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		PerPage:     perPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}
