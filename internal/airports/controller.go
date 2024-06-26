package airports

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"gitlab.strale.io/go-travel/internal/utils/handler/dto"
)

type airportService interface {
	ListAirports(ctx context.Context, pagination utils.Pagination) ([]database.Airport, error)
	ListAirportsInCity(ctx context.Context, cityID int64, pagination utils.Pagination) ([]database.Airport, error)
	FindByID(ctx context.Context, id int64) (database.Airport, error)
	SaveNewAirport(ctx context.Context, airport database.Airport) (database.Airport, error)
	UpdateAirport(ctx context.Context, airport database.Airport) error
	DeleteAirport(ctx context.Context, id int64) error
	ImportAirports(ctx context.Context, content []byte) error
}

type airportController struct {
	airportSrvc airportService
}

func NewAirportController(airportSrvc airportService) *airportController {
	return &airportController{
		airportSrvc: airportSrvc,
	}
}

func (ac *airportController) RegisterHandlers(airportPrefixed *mux.Router, cityPrefixed *mux.Router) {
	airportPrefixed.Path("").Methods(http.MethodGet).HandlerFunc(ac.ListAllAirports)
	airportPrefixed.Path("").Methods(http.MethodPost).HandlerFunc(ac.SaveNewAirport)
	airportPrefixed.Path("").Methods(http.MethodPatch).HandlerFunc(ac.ImportAirports)

	airportPrefixed.Path("/{id}").Methods(http.MethodGet).HandlerFunc(ac.GetAirportByID)
	airportPrefixed.Path("/{id}").Methods(http.MethodPut).HandlerFunc(ac.UpdateAirport)
	airportPrefixed.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(ac.DeleteAirport)

	cityPrefixed.Path("/{id}/airports").Methods(http.MethodGet).HandlerFunc(ac.ListAirportsInCity)
}

func (ac *airportController) listAirports(
	w http.ResponseWriter,
	r *http.Request,
	getter func(pagination utils.Pagination) ([]database.Airport, error),
) {
	pagination, ok := utils.PaginationFromRequest(w, r)
	if !ok {
		return
	}
	airports, err := getter(pagination)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.AirportsToDtos(airports))
}

func (ac *airportController) ListAllAirports(w http.ResponseWriter, r *http.Request) {
	ac.listAirports(
		w,
		r,
		func(pagination utils.Pagination) ([]database.Airport, error) {
			return ac.airportSrvc.ListAirports(r.Context(), pagination)
		},
	)
}

func (ac *airportController) ListAirportsInCity(w http.ResponseWriter, r *http.Request) {
	cityID, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	ac.listAirports(
		w,
		r,
		func(pagination utils.Pagination) ([]database.Airport, error) {
			return ac.airportSrvc.ListAirportsInCity(r.Context(), cityID.Val(), pagination)
		},
	)
}

func (ac *airportController) GetAirportByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	airport, err := ac.airportSrvc.FindByID(r.Context(), id.Val())
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.AirportToDto(airport))
}

func (ac *airportController) SaveNewAirport(w http.ResponseWriter, r *http.Request) {
	var payload dto.SaveAirportDto
	err := handler.GetBody(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	airport, err := ac.airportSrvc.SaveNewAirport(r.Context(), database.Airport{
		Name:   payload.Name,
		CityID: payload.CityID,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusCreated, dto.AirportToDto(airport))
}

func (ac *airportController) UpdateAirport(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	var payload dto.SaveAirportDto
	err = handler.GetBody(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = ac.airportSrvc.UpdateAirport(r.Context(), database.Airport{
		ID:     id.Val(),
		Name:   payload.Name,
		CityID: payload.CityID,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (ac *airportController) DeleteAirport(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = ac.airportSrvc.DeleteAirport(r.Context(), id.Val())
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (ac *airportController) ImportAirports(w http.ResponseWriter, r *http.Request) {
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
	go ac.airportSrvc.ImportAirports(context.Background(), bytes)
	handler.Respond(w, http.StatusOK, nil)
}
