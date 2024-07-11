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

type iPathFindingService interface {
	FindCheapestPath(ctx context.Context, startCityID, finishCityID int64) ([]*database.Route, int64, error)
}

type routeController struct {
	routeSrvc iRouteService
	pfSrvc    iPathFindingService
}

func NewRouteController(routeSrvc iRouteService, pfSrvc iPathFindingService) *routeController {
	return &routeController{
		routeSrvc: routeSrvc,
		pfSrvc:    pfSrvc,
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
	handler.OptionsAllowedMethods(input.RoutesRouter, "", http.MethodGet, http.MethodPost)

	input.RoutesRouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(rc.findByID)
	input.RoutesRouter.Path("/{id}").Methods(http.MethodPut).HandlerFunc(rc.update)
	input.RoutesRouter.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(rc.delete)
	handler.OptionsAllowedMethods(input.RoutesRouter, "/{id}", http.MethodGet, http.MethodPut, http.MethodDelete)

	input.AirportRouter.Path("/{id}/routes").Methods(http.MethodGet).HandlerFunc(rc.listRoutesForAirport)
	handler.OptionsAllowedMethods(input.AirportRouter, "/{id}/routes", http.MethodGet)

	input.CityRouter.Path("/{id}/routes").Methods(http.MethodGet).HandlerFunc(rc.listRoutesForCity)
	handler.OptionsAllowedMethods(input.CityRouter, "/{id}/routes", http.MethodGet)

	input.V1Router.Path("/cheapest-route").Methods(http.MethodGet).HandlerFunc(rc.cheapestPath)
	handler.OptionsAllowedMethods(input.V1Router, "/cheapest-route", http.MethodGet)
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
	cityID, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	rc.doList(w, r, func(pagination utils.Pagination) ([]database.Route, error) {
		return rc.routeSrvc.ListRoutesForCity(
			r.Context(),
			cityID.Val(),
			getDirection(r),
			pagination,
		)
	})
}

func (rc *routeController) listRoutesForAirport(w http.ResponseWriter, r *http.Request) {
	airportID, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	rc.doList(w, r, func(pagination utils.Pagination) ([]database.Route, error) {
		return rc.routeSrvc.ListRoutesForAirport(
			r.Context(),
			airportID.Val(),
			getDirection(r),
			pagination,
		)
	})
}

func (rc *routeController) findByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	route, err := rc.routeSrvc.RouteByID(r.Context(), id.Val())
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
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
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
	err = rc.routeSrvc.UpdateRoutePrice(r.Context(), id.Val(), payload.Price)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (rc *routeController) delete(w http.ResponseWriter, r *http.Request) {
	id, err := handler.Path(r, handler.Int64, "id", handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	err = rc.routeSrvc.DeleteRoute(r.Context(), id.Val())
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	handler.Respond(w, http.StatusOK, nil)
}

func (rc *routeController) cheapestPath(w http.ResponseWriter, r *http.Request) {
	beginID, err := handler.Query(r, handler.Int64, "begin", true, 0, handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	endID, err := handler.Query(r, handler.Int64, "end", true, 0, handler.IsPositive)
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	path, cheapestPrice, err := rc.pfSrvc.FindCheapestPath(r.Context(), beginID.Val(), endID.Val())
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	payload, err := dto.CompileCheapestPath(path, cheapestPrice).Encode()
	if err != nil {
		handler.ResolveErrorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}
