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
func PaginationFromRequest(w http.ResponseWriter, r *http.Request, resp *handler.Responder) (Pagination, bool) {
	page, err := handler.Query(r, handler.Int, "page", false, 0, handler.IsZeroOrPositive)
	if err != nil {
		resp.ResolveErrorResponse(w, err)
		return Pagination{}, false
	}
	pageSize, err := handler.Query(r, handler.Int, "page-size", false, 10, handler.IsZeroOrPositive)
	if err != nil {
		resp.ResolveErrorResponse(w, err)
		return Pagination{}, false
	}
	return Pagination{
		Limit:  pageSize.Val(),
		Offset: page.Val() * pageSize.Val(),
	}, true
}
