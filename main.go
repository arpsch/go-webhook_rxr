package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	api_http "github.com/arpsch/go-webhook_rxr/api/http"
	"github.com/arpsch/go-webhook_rxr/imemc"
	"github.com/arpsch/go-webhook_rxr/receiver"
)

var (
	port = ":9999"
)

func setupServer() (*httprouter.Router, error) {

	// set up receiver instance
	whRxr := receiver.NewReceiver()

	// construct inmemory cache
	imc := imemc.NewLogCache()

	// set up the go routine to handle cacche evictions
	go imemc.HandleWebhookEvents(imc)

	// set up api handlers for image collector
	whHandler := api_http.NewReceiverHandlers(whRxr, imc)
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
