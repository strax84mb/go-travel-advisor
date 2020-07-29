package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
	h "gitlab.strale.io/go-travel/handlers"
)

// Hello a handler for /hello/:name endpoint
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	content := fmt.Sprintf("Hello to you too %s\n", ps.ByName("name"))
	w.Write([]byte(content))
}

func main() {
	db.InitDb()

	router := httprouter.New()
	router.GET("/hello/:name", Hello)
	// user endpoints
	router.POST("/user/signup", h.SignupUser)
	router.POST("/user/login", h.LoginUser)
	// city endpoints
	router.GET("/city", h.GetAllCities)
	router.POST("/city", h.AddCity)
	router.GET("/city/:id", h.GetCity)
	router.PUT("/city/:id", h.UpdateCity)
	router.DELETE("/city/:id", h.DeleteCity)
	router.POST("/city/import", h.ImportCities)
	// comment endpoints
	router.POST("/comment", h.PostComment)
	router.PUT("/comment/:id", h.UpdateComment)
	router.DELETE("/comment/:id", h.DeleteComment)
	// airport endpoints
	router.GET("/airport", h.ListAirports)
	router.POST("/airport", h.AddAirport)
	router.GET("/airport/:id", h.GetAirport)
	router.PUT("/airport/:id", h.UpdateAirport)
	router.DELETE("/airport", h.DeleteAirport)
	router.POST("/airport/import", h.ImportAirports)
	// route endpoints

	log.Fatal(http.ListenAndServe(":8080", router))
}
