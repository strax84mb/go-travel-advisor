package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/airports"
	airportRepo "gitlab.strale.io/go-travel/internal/airports/repository"
	cities "gitlab.strale.io/go-travel/internal/cities"
	cityRepo "gitlab.strale.io/go-travel/internal/cities/repository"
	"gitlab.strale.io/go-travel/internal/comments"
	commentRepo "gitlab.strale.io/go-travel/internal/comments/repository"
	"gitlab.strale.io/go-travel/internal/config"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/middleware"
	"gitlab.strale.io/go-travel/internal/routes"
	routeRepo "gitlab.strale.io/go-travel/internal/routes/repository"
	"gitlab.strale.io/go-travel/internal/users"
	userRepo "gitlab.strale.io/go-travel/internal/users/repository"
	"gitlab.strale.io/go-travel/internal/utils"
)

// Hello a handler for /hello/:name endpoint
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	content := fmt.Sprintf("Hello to you too %s\n", ps.ByName("name"))
	w.Write([]byte(content))
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.AddHook(&utils.ContextLoggerHook{})

	conf, err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatal("failed to load configuration ", err.Error())
	}

	db, err := database.ConnectToDatabase(conf.DB)
	if err != nil {
		log.Fatal("failed to connect to database ", err.Error())
	}

	cityRepository := cityRepo.NewCityRepository(&db)
	airportRepository := airportRepo.NewAirportRepository(&db)
	userRepository := userRepo.NewUserRepository(&db)
	commentRepository := commentRepo.NewCommentRepository(&db)
	routeRepository := routeRepo.NewRouteRepository(&db)

	cityService := cities.NewCityService(cityRepository, airportRepository)
	airportService := airports.NewAirportService(airportRepository, cityRepository)
	commentService := comments.NewCommentService(
		commentRepository,
		cityRepository,
		userRepository,
	)
	routesService := routes.NewRouteService(routeRepository, airportRepository)
	securityService, err := users.NewSecurityService(conf.Security.RSAKey)
	if err != nil {
		log.Fatal("failed to initialize security ", err.Error())
	}
	userService := users.NewUserService(userRepository, securityService)
	pathFindingService := routes.NewPathFindingService(
		airportRepository,
		cityRepository,
		routeRepository,
	)

	cityController := cities.NewCityController(cityService)
	airportController := airports.NewAirportController(airportService)
	commentsController := comments.NewCommentController(commentService)
	routesController := routes.NewRouteController(routesService, pathFindingService)
	usersController := users.NewUserController(userService)

	jwtMiddleware := middleware.NewVerifyJWTMiddleware(securityService)
	corsMiddleware := middleware.NewCorsMiddleware(conf.Client.Origin)

	r := mux.NewRouter()
	r.Use(middleware.RequestIDMiddleware)
	r.Use(jwtMiddleware.Middleware)
	r.Use(corsMiddleware.Middleware)

	v1Router := r.PathPrefix("/v1").Subrouter()
	cityPrefixed := v1Router.PathPrefix("/cities").Subrouter()
	commentPrefixed := v1Router.PathPrefix("/comments").Subrouter()
	userPrefixed := v1Router.PathPrefix("/users").Subrouter()
	routePrefixed := v1Router.PathPrefix("/routes").Subrouter()
	airportPrefixed := v1Router.PathPrefix("/airports").Subrouter()

	cityController.RegisterHandlers(cityPrefixed)
	airportController.RegisterHandlers(airportPrefixed, cityPrefixed)
	commentsController.RegisterHandlers(comments.RegisterHandlersInput{
		V1Prefixed:       v1Router,
		CityPrefixed:     cityPrefixed,
		UsersPrefixed:    userPrefixed,
		CommentsPrefixed: commentPrefixed,
	})
	routesController.RegisterHandlers(routes.RegisterHandlersInput{
		V1Router:      v1Router,
		RoutesRouter:  routePrefixed,
		CityRouter:    cityPrefixed,
		AirportRouter: airportPrefixed,
	})
	usersController.RegisterHandlers(v1Router, userPrefixed)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", conf.Server.Address, conf.Server.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Info("Serving on port ", conf.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
