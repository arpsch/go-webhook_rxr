package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	api_http "github.com/arpsch/go-webhook_rxr/api/http"
	"github.com/arpsch/go-webhook_rxr/client"
	"github.com/arpsch/go-webhook_rxr/imemc"
)

var (
	port = ":9999"
)

func setupServer() (*httprouter.Router, error) {

	// construct inmemory cache
	imc := imemc.NewLogCache()

	// create client object
	client := client.NewClient()

	// set up api handlers for webhook receiver
	whHandler := api_http.NewReceiverHandlers(imc, client)

	// set up the go routine to handle cacche evictions
	go api_http.HandleWebhookEvents(whHandler)

	// set up the routes
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
