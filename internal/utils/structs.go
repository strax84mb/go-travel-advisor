package utils

import (
	"net/http"

	"gitlab.strale.io/go-travel/internal/utils/handler"
)

// Pagination object
type Pagination struct {
	Limit  int
	Offset int
}

func PaginationFrom(page, pageSize int) Pagination {
	return Pagination{
		Limit:  pageSize,
		Offset: page * pageSize,
	}
}

// Get pagination query parameters from request
//
//	pagination, ok := cc.getPagination(w, r)
//
// Parameters:
//   - w --> http.ResponseWriter used to write error response
//   - r --> *http.Request from which to extract
//     query parameters from
//
// Returns:
//   - pagination --> pagination object
//   - ok --> no error happened and operation was successful
func PaginationFromRequest(w http.ResponseWriter, r *http.Request) (Pagination, bool) {
	page, err := handler.QueryAsInt(r, "page", false, 0, handler.IntMustBeZeroOrPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return Pagination{}, false
	}
	pageSize, err := handler.QueryAsInt(r, "page-size", false, 10, handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return Pagination{}, false
	}
	return Pagination{
		Limit:  pageSize,
		Offset: page * pageSize,
	}, true
}
