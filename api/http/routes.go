package http

import (
	"github.com/arpsch/go-webhook_rxr/logger"
	"github.com/julienschmidt/httprouter"
)

// SetupRoutes sets up routes for the WH reveiver
func (rh *receiverHandlers) SetupRoutes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc("GET", "/healthz", logger.Log(rh.HealthzHandler))
	router.HandlerFunc("POST", "/log", logger.Log(rh.HookHandler))

	return router
}
