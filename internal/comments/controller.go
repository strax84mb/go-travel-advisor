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
	ListComments(ctx context.Context, input utils.Pagination) ([]database.Comment, error)
	FindByID(ctx context.Context, id int64) (*database.Comment, error)
	ListCommentsForCity(ctx context.Context, cityID int64, pagination utils.Pagination) ([]database.Comment, error)
	ListCommentsForUser(ctx context.Context, userID int64, pagination utils.Pagination) ([]database.Comment, error)
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

type RegisterHandlersInput struct {
	V1Prefixed       *mux.Router
	CityPrefixed     *mux.Router
	UsersPrefixed    *mux.Router
	CommentsPrefixed *mux.Router
}

func (cc *commentController) RegisterHandlers(input RegisterHandlersInput) {
	input.V1Prefixed.Path("/me/comments").Methods(http.MethodGet).HandlerFunc(cc.listCommentsForMe)

	input.CityPrefixed.Path("/cities/{id}/comments").Methods(http.MethodGet).HandlerFunc(cc.listCommentsForCity)

	input.UsersPrefixed.Path("/users/{id}/comments").Methods(http.MethodGet).HandlerFunc(cc.listCommentsForUser)

	input.CommentsPrefixed.Path("/comments").Methods(http.MethodGet).HandlerFunc(cc.listComments)
	input.CommentsPrefixed.Path("/comments").Methods(http.MethodPost).HandlerFunc(cc.saveNewComment)

	input.CommentsPrefixed.Path("/comments/{id}").Methods(http.MethodGet).HandlerFunc(cc.getCommentByID)
	input.CommentsPrefixed.Path("/comments/{id}").Methods(http.MethodPut).HandlerFunc(cc.updateComment)
	input.CommentsPrefixed.Path("/comments/{id}").Methods(http.MethodDelete).HandlerFunc(cc.deleteComment)

	input.CommentsPrefixed.Path("/comments/{id}/force").Methods(http.MethodDelete).HandlerFunc(cc.forceDeleteComment)
}

func (cc *commentController) listComments(w http.ResponseWriter, r *http.Request) {
	pagination, ok := utils.PaginationFromRequest(w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	comments, err := cc.commentSrvc.ListComments(ctx, pagination)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) getCommentByID(w http.ResponseWriter, r *http.Request) {
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
	handler.RespondFF(w, http.StatusOK, dto.CommentToDto(*comment))
}

func (cc *commentController) listCommentsForMe(w http.ResponseWriter, r *http.Request) {
	pagination, ok := utils.PaginationFromRequest(w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	userID, _, ok := utils.GetJWTData(ctx)
	if !ok {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("user not logged in"))
		return
	}
	comments, err := cc.commentSrvc.ListCommentsForUser(ctx, userID, pagination)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) listCommentsForUser(w http.ResponseWriter, r *http.Request) {
	pagination, ok := utils.PaginationFromRequest(w, r)
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
	comments, err := cc.commentSrvc.ListCommentsForUser(ctx, userID, pagination)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) listCommentsForCity(w http.ResponseWriter, r *http.Request) {
	pagination, ok := utils.PaginationFromRequest(w, r)
	if !ok {
		return
	}
	cityID, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	comments, err := cc.commentSrvc.ListCommentsForCity(r.Context(), cityID, pagination)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusOK, dto.CommentsToDtos(comments))
}

func (cc *commentController) saveNewComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, ok := utils.GetJWTData(ctx)
	if !ok {
		handler.ResolveErrorResponse(w, users.NewErrUnauthorizedWithCause("must be logged in"))
		return
	}
	var payload dto.SaveCommentDto
	err := handler.GetBodyFF(r, &payload)
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
	handler.RespondFF(w, http.StatusCreated, dto.CommentToDto(*comment))
}

func (cc *commentController) updateComment(w http.ResponseWriter, r *http.Request) {
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
	err = handler.GetBodyFF(r, &payload)
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

func (cc *commentController) deleteComment(w http.ResponseWriter, r *http.Request) {
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

func (cc *commentController) forceDeleteComment(w http.ResponseWriter, r *http.Request) {
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
