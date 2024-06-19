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
}

func NewUserController(userSrvc iUserService) *userController {
	return &userController{
		userSrvc: userSrvc,
	}
}

func (uc *userController) RegisterHandlers(v1Router *mux.Router, userRouter *mux.Router) {
	v1Router.Path("/login").Methods(http.MethodPost).HandlerFunc(uc.Login)

	userRouter.Path("").Queries("username", "{username:.+}").Methods(http.MethodGet).HandlerFunc(uc.GetUserByUsername)
	userRouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(uc.GetUserById)
}

func (uc *userController) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin dto.UserLoginDto
	err := handler.GetBody(r, &userLogin)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	token, err := uc.userSrvc.Login(r.Context(), userLogin.Username, userLogin.Password)
	if err != nil {
		handler.ResolveErrorResponse(w, handler.NewErrUnauthorizedWithCause("invalid username or password"))
		return
	}
	payload := dto.LoginTokenDto{
		Token: token,
	}
	handler.Respond(w, http.StatusOK, &payload)
}

func (uc *userController) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		handler.ResolveErrorResponse(w, handler.NewErrBadRequest("username missing"))
		return
	}
	user, err := uc.userSrvc.GetByUsername(r.Context(), username)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.UserToDto(user))
}

func (uc *userController) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64FromString, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	user, err := uc.userSrvc.GetByID(r.Context(), id.Val())
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.UserToDto(user))
}
