package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	api_http "github.com/arpsch/go-webhook_rxr/api/http"
	"github.com/arpsch/go-webhook_rxr/receiver"
)

var (
	port = ":9999"
)

func setupServer() (*httprouter.Router, error) {

	// set up receiver instance
	whRxr := receiver.NewReceiver()

	// set up api handlers for image collector
	whHandler := api_http.NewReceiverHandlers(whRxr)
	routes := whHandler.SetupRoutes()

	return routes, nil
}

func main() {
	router, err := setupServer()
	if err != nil {
		log.Printf(" failed to set up routes, exiting")
		return
	}

	log.Printf("Server is listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
