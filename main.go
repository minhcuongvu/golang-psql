package main

import (
	"counter-service/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
