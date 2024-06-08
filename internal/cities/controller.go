package cities

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"gitlab.strale.io/go-travel/internal/utils/handler/dto"
)

type iCityService interface {
	ListCities(ctx context.Context, offset, limit int) ([]database.City, error)
	FindByID(ctx context.Context, id int64) (database.City, error)
	SaveNewCity(ctx context.Context, name string) (database.City, error)
	UpdateCity(ctx context.Context, id int64, name string) error
	DeleteCity(ctx context.Context, id int64) error
	ImportCities(ctx context.Context, content []byte) error
}

type cityController struct {
	citySrvc iCityService
}

func NewCityController(citySrvc iCityService) *cityController {
	return &cityController{
		citySrvc: citySrvc,
	}
}

func (cc *cityController) RegisterHandlers(r *mux.Router) {
	r.Path("").Methods(http.MethodGet).HandlerFunc(cc.ListAllCities)
	r.Path("").Methods(http.MethodPost).HandlerFunc(cc.SaveNewCity)
	r.Path("").Methods(http.MethodPatch).HandlerFunc(cc.ImportCities)

	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(cc.GetCityByID)
	r.Path("/{id}").Methods(http.MethodPut).HandlerFunc(cc.UpdateCity)
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(cc.DeleteCity)
}

func (cc *cityController) ListAllCities(w http.ResponseWriter, r *http.Request) {
	page, err := handler.QueryAsInt(r, "page", false, 0, handler.IntMustBeZeroOrPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	pageSize, err := handler.QueryAsInt(r, "page-size", false, 10, handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	ctx := r.Context()
	cities, err := cc.citySrvc.ListCities(ctx, page*pageSize, pageSize)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusOK, dto.CitiesToDtos(cities))
}

func (cc *cityController) GetCityByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	ctx := r.Context()
	city, err := cc.citySrvc.FindByID(ctx, id)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusOK, dto.CityToDto(city))
}

func (cc *cityController) SaveNewCity(w http.ResponseWriter, r *http.Request) {
	var payload dto.SaveCityDto
	err := handler.GetBodyFF(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	city, err := cc.citySrvc.SaveNewCity(r.Context(), payload.Name)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.RespondFF(w, http.StatusCreated, dto.CityToDto(city))
}

func (cc *cityController) UpdateCity(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	var payload dto.SaveCityDto
	err = handler.GetBodyFF(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	ctx := r.Context()
	err = cc.citySrvc.UpdateCity(ctx, id, payload.Name)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (cc *cityController) DeleteCity(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = cc.citySrvc.DeleteCity(r.Context(), id)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (cc *cityController) ImportCities(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		handler.ResolveErrorResponse(w, handler.NewErrBadRequest(
			err.Error(),
		))
		return
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		handler.ResolveErrorResponse(w, handler.NewErrBadRequest(
			fmt.Sprintf("could not read contents of uploaded file: %s", err.Error()),
		))
		return
	}
	go cc.citySrvc.ImportCities(context.Background(), bytes)
	handler.Respond(w, http.StatusOK, nil)
}
