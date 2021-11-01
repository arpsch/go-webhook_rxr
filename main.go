package main

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	api_http "github.com/arpsch/go-webhook_rxr/api/http"
	"github.com/arpsch/go-webhook_rxr/client"
	"github.com/arpsch/go-webhook_rxr/imemc"
	log "github.com/arpsch/go-webhook_rxr/logger"
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
	l := log.Logger{}

	router, err := setupServer()
	if err != nil {
		l.Log(log.PANIC, "failed to set up routes, exiting")
		return
	}

	srv := &http.Server{
		Addr:    port,
		Handler: router,

		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	l.Log(log.INFO, "server is listening on %s", port)
	l.Log(log.PANIC, srv.ListenAndServe().Error())
}
