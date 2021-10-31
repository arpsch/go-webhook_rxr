package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arpsch/go-webhook_rxr/client"
	"github.com/arpsch/go-webhook_rxr/imemc"
	"github.com/arpsch/go-webhook_rxr/model"
	"github.com/pkg/errors"
)

// TODO: Read from environment variable
const (
	BATCHSIZE     = 3 // number of items in cache
	BATCHINTERVAL = 30 * time.Second
)

type receiverHandlers struct {
	cache  *imemc.Cache
	client *client.Client

	// signal for go routine
	BatchSizeChan chan struct{}
}

// NewReceiverHandlers constructor for Receiver
func NewReceiverHandlers(cache *imemc.Cache, client *client.Client) *receiverHandlers {
	return &receiverHandlers{
		cache:  cache,
		client: client,

		BatchSizeChan: make(chan struct{}),
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

	log.Printf("cache write at %v \n", time.Now())
	batchSize, err := rh.cache.Write(ctx, l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if batchSize >= BATCHSIZE {
		rh.BatchSizeChan <- struct{}{}
	}

	w.WriteHeader(http.StatusAccepted)
}

// HandleWebhookEvents handle the webhooks based on requirents
func HandleWebhookEvents(rh *receiverHandlers) {
	for {
		t := time.NewTimer(BATCHINTERVAL)
		select {
		case <-rh.BatchSizeChan:
			log.Printf("cache batch size threshold crossed")
			logs, err := rh.cache.Evict()
			if err == nil {
				if err := client.RetryTimeout(logs, rh.client.HandleHook); err != nil {
					log.Printf("failed to release events to upstream server: %v", err)
				}
			}
		case <-t.C:
			log.Printf("cache batch interval crossed")
			logs, err := rh.cache.Evict()
			if err == nil {
				log.Printf("current batch size: %d", len(logs))
				// if empty cache do nothing
				if len(logs) <= 0 {
					continue
				}
				if err := client.RetryTimeout(logs, rh.client.HandleHook); err != nil {
					log.Printf("failed to release events to upstream server: %v", err)
				}
			}
		}
	}
}
