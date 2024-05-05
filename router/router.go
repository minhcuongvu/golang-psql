package router

import (
	"counter-service/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/increment", handlers.IncrementAndFetch).Methods("GET", "OPTIONS")
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	return router
}
