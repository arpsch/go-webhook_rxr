package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arpsch/go-webhook_rxr/model"
	"github.com/arpsch/go-webhook_rxr/receiver"
	"github.com/pkg/errors"
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

func parseLog(ctx context.Context, r *http.Request) (*model.Log, error) {
	log := model.Log{}

	//decode body
	err := json.NewDecoder(r.Body).Decode(&log)

	switch {
	case err != nil:
		return nil, errors.Wrap(err, "failed to decode json")

		/*
			TODO: struct comparison is not working for reason :(
				to check. Working in playground.
			case log == model.Log{}:
			return nil, errors.New("empty request body")

		*/
	}
	return &log, nil
}

// HealthzHandler health update handler funcation
func (rh *receiverHandlers) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	// just return OK status and message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HookHandler webhook handler function
func (rh *receiverHandlers) HookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	l, err := parseLog(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", l)

	w.WriteHeader(http.StatusAccepted)
}
