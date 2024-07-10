package users

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"gitlab.strale.io/go-travel/internal/utils/handler/dto"
)

type iUserService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetByUsername(ctx context.Context, username string) (*database.User, error)
	GetByID(ctx context.Context, id int64) (*database.User, error)
}

type userController struct {
	userSrvc iUserService
	r        *handler.Responder
}

func NewUserController(userSrvc iUserService, r *handler.Responder) *userController {
	return &userController{
		userSrvc: userSrvc,
		r:        r,
	}
}

func (uc *userController) RegisterHandlers(v1Router *mux.Router, userRouter *mux.Router, c *handler.Cors) {
	v1Router.Path("/login").Methods(http.MethodPost).HandlerFunc(uc.Login)
	c.Options(v1Router, "/login", http.MethodPost)

	userRouter.Path("").Queries("username", "{username:.+}").Methods(http.MethodGet).HandlerFunc(uc.GetUserByUsername)
	c.Options(userRouter, "", http.MethodGet)
	userRouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(uc.GetUserById)
	c.Options(userRouter, "/{id}", http.MethodGet)
}

func (uc *userController) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin dto.UserLoginDto
	err := handler.GetBody(r, &userLogin)
	if err != nil {
		uc.r.ResolveErrorResponse(w, err)
		return
	}
	token, err := uc.userSrvc.Login(r.Context(), userLogin.Username, userLogin.Password)
	if err != nil {
		uc.r.ResolveErrorResponse(w, handler.NewErrUnauthorizedWithCause("invalid username or password"))
		return
	}
	payload := dto.LoginTokenDto{
		Token: token,
	}
	uc.r.Respond(w, http.StatusOK, &payload)
}

func (uc *userController) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		uc.r.ResolveErrorResponse(w, handler.NewErrBadRequest("username missing"))
		return
	}
	user, err := uc.userSrvc.GetByUsername(r.Context(), username)
	if err != nil {
		uc.r.ResolveErrorResponse(w, err)
		return
	}
	uc.r.Respond(w, http.StatusOK, dto.UserToDto(user))
}

func (uc *userController) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		uc.r.ResolveErrorResponse(w, err)
		return
	}
	user, err := uc.userSrvc.GetByID(r.Context(), id.Val())
	if err != nil {
		uc.r.ResolveErrorResponse(w, err)
		return
	}
	uc.r.Respond(w, http.StatusOK, dto.UserToDto(user))
}
