package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"gitlab.strale.io/go-travel/internal/utils/handler/dto"
)

type iRouteService interface {
	ListRoutes(ctx context.Context, pagination utils.Pagination) ([]database.Route, error)
	ListRoutesForCity(ctx context.Context, cityID int64, direction Direction, pagination utils.Pagination) ([]database.Route, error)
	ListRoutesForAirport(ctx context.Context, airportID int64, direction Direction, pagination utils.Pagination) ([]database.Route, error)
	RouteByID(ctx context.Context, id int64) (*database.Route, error)
	SaveNew(ctx context.Context, route database.Route) (*database.Route, error)
	UpdateRoutePrice(ctx context.Context, id int64, price float32) error
	DeleteRoute(ctx context.Context, id int64) error
}

type routeController struct {
	routeSrvc iRouteService
}

func NewRouteController(routeSrvc iRouteService) *routeController {
	return &routeController{
		routeSrvc: routeSrvc,
	}
}

type RegisterHandlersInput struct {
	V1Router      *mux.Router
	RoutesRouter  *mux.Router
	CityRouter    *mux.Router
	AirportRouter *mux.Router
}

func (rc *routeController) RegisterHandlers(input RegisterHandlersInput) {
	input.RoutesRouter.Path("").Methods(http.MethodGet).HandlerFunc(rc.listRoutes)
	input.RoutesRouter.Path("").Methods(http.MethodPost).HandlerFunc(rc.saveNewRoute)

	input.RoutesRouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(rc.findByID)
	input.RoutesRouter.Path("/{id}").Methods(http.MethodPut).HandlerFunc(rc.update)
	input.RoutesRouter.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(rc.delete)

	input.AirportRouter.Path("/{id}/routes").Methods(http.MethodGet).HandlerFunc(rc.listRoutesForAirport)

	input.CityRouter.Path("/{id}/routes").Methods(http.MethodGet).HandlerFunc(rc.listRoutesForCity)

	// TODO implement this
	// /cheapest-route?begin=?&end=?
}

func (rc *routeController) doList(
	w http.ResponseWriter,
	r *http.Request,
	getter func(utils.Pagination) ([]database.Route, error),
) {
	pagination, ok := utils.PaginationFromRequest(w, r)
	if !ok {
		return
	}
	list, err := getter(pagination)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.RoutesToDtos(list))
}

func (rc *routeController) listRoutes(w http.ResponseWriter, r *http.Request) {
	rc.doList(w, r, func(pagination utils.Pagination) ([]database.Route, error) {
		return rc.routeSrvc.ListRoutes(r.Context(), pagination)
	})
}

func getDirection(r *http.Request) Direction {
	switch r.URL.Query().Get("direction") {
	case "ALL":
		return DIRECTION_ALL
	case "INCOMMING":
		return DIRECTION_INCOMMING
	case "OUTGOING":
		return DIRECTION_OUTGOING
	default:
		return DIRECTION_ALL
	}
}

func (rc *routeController) listRoutesForCity(w http.ResponseWriter, r *http.Request) {
	cityID, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	rc.doList(w, r, func(pagination utils.Pagination) ([]database.Route, error) {
		return rc.routeSrvc.ListRoutesForCity(
			r.Context(),
			cityID,
			getDirection(r),
			pagination,
		)
	})
}

func (rc *routeController) listRoutesForAirport(w http.ResponseWriter, r *http.Request) {
	airportID, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	rc.doList(w, r, func(pagination utils.Pagination) ([]database.Route, error) {
		return rc.routeSrvc.ListRoutesForAirport(
			r.Context(),
			airportID,
			getDirection(r),
			pagination,
		)
	})
}

func (rc *routeController) findByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	route, err := rc.routeSrvc.RouteByID(r.Context(), id)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusCreated, dto.RouteToDto(*route))
}

func (rc *routeController) saveNewRoute(w http.ResponseWriter, r *http.Request) {
	var payload dto.SaveRouteDto
	err := handler.GetBody(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	route, err := rc.routeSrvc.SaveNew(r.Context(), database.Route{
		SourceID:      payload.SourceID,
		DestinationID: payload.DestinationID,
		Price:         payload.Price,
	})
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, dto.RouteToDto(*route))
}

func (rc *routeController) update(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	var payload dto.UpdateRoutePrice
	err = handler.GetBody(r, &payload)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = rc.routeSrvc.UpdateRoutePrice(r.Context(), id, payload.Price)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (rc *routeController) delete(w http.ResponseWriter, r *http.Request) {
	id, err := handler.PathAsInt64(r, "id", handler.IntMustBePositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = rc.routeSrvc.DeleteRoute(r.Context(), id)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}
