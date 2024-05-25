package comments

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/users"
	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"gitlab.strale.io/go-travel/internal/utils/handler/dto"
)

type iCommentService interface {
	ListComments(ctx context.Context, input ListCommentsInput) ([]database.Comment, error)
	FindByID(ctx context.Context, id int64) (*database.Comment, error)
	ListCommentsForCity(ctx context.Context, cityID int64, pagination ListCommentsInput) ([]database.Comment, error)
	ListCommentsForUser(ctx context.Context, userID int64, pagination ListCommentsInput) ([]database.Comment, error)
	SaveComment(ctx context.Context, comment database.Comment) (*database.Comment, error)
	UpdateText(ctx context.Context, commentID, requestorID int64, text string) error
	DeleteByID(ctx context.Context, commentID, requestorID int64, force bool) error
}

type commentController struct {
	commentSrvc iCommentService
}

func NewCommentController(commentSrvc iCommentService) *commentController {
	return &commentController{
		commentSrvc: commentSrvc,
	}
}

func (cc *commentController) RegisterHandlers(
	v1Prefixed *mux.Router,
	cityPrefixed *mux.Router,
	usersPrefixed *mux.Router,
	commentsPrefixed *mux.Router,
) {
	v1Prefixed.Path("/me/comments").Methods(http.MethodGet).HandlerFunc(cc.ListCommentsForMe)

	cityPrefixed.Path("/cities/{id}/comments").Methods(http.MethodGet).HandlerFunc(cc.ListCommentsForCity)

	usersPrefixed.Path("/users/{id}/comments").Methods(http.MethodGet).HandlerFunc(cc.ListCommentsForUser)

	commentsPrefixed.Path("/comments").Methods(http.MethodGet).HandlerFunc(cc.ListComments)
	commentsPrefixed.Path("/comments").Methods(http.MethodPost).HandlerFunc(cc.SaveNewComment)

	commentsPrefixed.Path("/comments/{id}").Methods(http.MethodGet).HandlerFunc(cc.GetCommentByID)
	commentsPrefixed.Path("/comments/{id}").Methods(http.MethodPut).HandlerFunc(cc.UpdateComment)
	commentsPrefixed.Path("/comments/{id}").Methods(http.MethodDelete).HandlerFunc(cc.DeleteComment)

	commentsPrefixed.Path("/comments/{id}/force").Methods(http.MethodDelete).HandlerFunc(cc.ForceDeleteComment)
}

// Get pagination query parameters from request
//
//	page, pageSize, ok := cc.getPagination(w, r)
//
// Parameters:
//   - w --> http.ResponseWriter used to write error response
//   - r --> *http.Request from which to extract
//     query parameters from
//
// Returns:
//   - page --> requested page
//   - pageSize --> size of requested page
//   - ok --> no error happened and operation was successful
func (cc *commentController) getPagination(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	page, err := handler.QueryAsInt(r, "page", false, 0, handler.IntMustBeZeroOrPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return 0, 0, false
	}
	pageSize, err := handler.QueryAsInt(r, "page-size", false, 10, handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return 0, 0, false
	}
	return page, pageSize, true
}

func (cc *commentController) ListComments(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := cc.getPagination(w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	comments, err := cc.commentSrvc.ListComments(ctx, ListCommentsInput{
		limit:  pageSize,
		offset: page * pageSize,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	comment, err := cc.commentSrvc.FindByID(r.Context(), id)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.CommentToDto(*comment))
}

func (cc *commentController) ListCommentsForMe(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := cc.getPagination(w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	userID, _, ok := utils.GetJWTData(ctx)
	if !ok {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("user not logged in"))
		return
	}
	comments, err := cc.commentSrvc.ListCommentsForUser(ctx, userID, ListCommentsInput{
		limit:  pageSize,
		offset: page * pageSize,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) ListCommentsForUser(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := cc.getPagination(w, r)
	if !ok {
		return
	}
	userID, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	ctx := r.Context()
	_, roles, ok := utils.GetJWTData(ctx)
	if !ok || !utils.HasRole(roles, "admin") {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("only admins allowed"))
		return
	}
	comments, err := cc.commentSrvc.ListCommentsForUser(ctx, userID, ListCommentsInput{
		limit:  pageSize,
		offset: page * pageSize,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) ListCommentsForCity(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := cc.getPagination(w, r)
	if !ok {
		return
	}
	cityID, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	comments, err := cc.commentSrvc.ListCommentsForCity(r.Context(), cityID, ListCommentsInput{
		limit:  pageSize,
		offset: page * pageSize,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) SaveNewComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, ok := utils.GetJWTData(ctx)
	if !ok {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("must be logged in"))
		return
	}
	var payload dto.SaveCommentDto
	err := handler.GetBody(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	comment, err := cc.commentSrvc.SaveComment(ctx, database.Comment{
		CityID:   payload.CityID,
		PosterID: userID,
		Text:     payload.Text,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusCreated, dto.CommentToDto(*comment))
}

func (cc *commentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, ok := utils.GetJWTData(ctx)
	if !ok {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("must be logged in"))
		return
	}
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	var payload dto.UpdateCommentDto
	err = handler.GetBody(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = cc.commentSrvc.UpdateText(r.Context(), id, userID, payload.Text)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (cc *commentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, ok := utils.GetJWTData(ctx)
	if !ok {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("must be logged in"))
		return
	}
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = cc.commentSrvc.DeleteByID(ctx, id, userID, false)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (cc *commentController) ForceDeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, roles, ok := utils.GetJWTData(ctx)
	if !ok || utils.HasRole(roles, "admin") {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("must be logged in"))
		return
	}
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = cc.commentSrvc.DeleteByID(ctx, id, 0, true)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}
