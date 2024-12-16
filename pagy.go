// Package pagy is a package for collecting pagination data from clients and returning a consistent paginated structure.
/*
pagy is not limited to any one framework or library, it works directly with the stdlib http.Request interface
so any project/framework that supports stdlib it will support.

pagy works by collecting pagination based query params from the http.Request and formatting it into an expected
and consistent structure to use throughout your service. A set of tools is available for getting the pagination data
as well as structures that are generic-based for responding with a consistent data contract.
*/
package pagy

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultSize = 10
)

// PaginationResponse response type.
type PaginationResponse[T any] struct {
	TotalCount int  `json:"total_count"`
	TotalPages int  `json:"total_pages"`
	Page       int  `json:"page"`
	Size       int  `json:"size"`
	HasMore    bool `json:"has_more"`
	Values     []T  `json:"values"`
}

// PaginatedResponse returns a hydrated pagination response model.
func PaginatedResponse[T any](count int, pq *PaginationQuery, list []T) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		TotalCount: count,
		TotalPages: GetTotalPages(count, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    GetHasMore(pq.GetPage(), count, pq.GetSize()),
		Values:     list,
	}
}

// DefaultPaginationResponse returns a default ( empty ) pagination response model.
func DefaultPaginationResponse[T any](pq *PaginationQuery) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		TotalCount: 0,
		TotalPages: GetTotalPages(0, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    GetHasMore(pq.GetPage(), 0, pq.GetSize()),
		Values:     make([]T, 0),
	}
}

// PaginationQuery query params.
type PaginationQuery struct {
	Size     int    `json:"size,omitempty"`
	Page     int    `json:"page,omitempty"`
	OrderBy  string `json:"orderBy,omitempty"`
	OrderDir string `json:"orderDir,omitempty"`
}

// SetSize converts the size value to an integer and sets it.
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

// SetPage converts the page value to an integer and sets it.
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// SetOrderBy sets the column order string with col-name + dir-identifier.
func (q *PaginationQuery) SetOrderBy(orderByQuery, orderDirQuery string) {
	if orderByQuery == "" {
		q.OrderBy = "id"
	} else {
		q.OrderBy = orderByQuery
	}

	if orderDirQuery == "" || strings.EqualFold(orderDirQuery, "asc") {
		q.OrderBy = fmt.Sprintf("%s %s", q.OrderBy, "ASC")
	} else {
		q.OrderBy = fmt.Sprintf("%s %s", q.OrderBy, "DESC")
	}
}

// GetOffset returns the column offset value.
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit returns the limit of columns to return.
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// GetOrderBy returns the column order value.
func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}

// GetPage returns the offset or ( page ).
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// GetSize returns the number of requested rows.
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

// GetQueryString returns an example of the query string passed in the http request.
func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

// GetPaginationFromRequest pagination query struct from.
func GetPaginationFromRequest(r *http.Request) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(r.URL.Query().Get("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(r.URL.Query().Get("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(r.URL.Query().Get("orderBy"), r.URL.Query().Get("orderDir"))

	return q, nil
}

// GetTotalPages returns the total pages count.
func GetTotalPages(totalCount, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

// GetHasMore returns the has more pages value.
func GetHasMore(currentPage, totalCount, pageSize int) bool {
	return currentPage < totalCount/pageSize
}
