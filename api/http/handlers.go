package http

import (
	"net/http"

	"github.com/arpsch/go-webhook_rxr/receiver"
)

type receiverHandlers struct {
	rxr receiver.WebhookReceiver
}

// NewReceiverHandlers constructor for Receiver
func NewReceiverHandlers(rxr receiver.WebhookReceiver) *receiverHandlers {
	return &receiverHandlers{
		rxr: rxr,
	}
}

// HealthzHandler health update handler funcation
func (rh *receiverHandlers) HealthzHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

// HookHandler webhook handler function
func (rh *receiverHandlers) HookHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusAccepted)
}
