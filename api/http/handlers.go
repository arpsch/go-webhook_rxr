package http

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/arpsch/go-webhook_rxr/client"
	"github.com/arpsch/go-webhook_rxr/imemc"
	log "github.com/arpsch/go-webhook_rxr/logger"
	"github.com/arpsch/go-webhook_rxr/model"
	"github.com/pkg/errors"
)

const (
	DEF_BATCH_SIZE     = 3
	DEF_BATCH_INTERVAL = 30 * time.Second
)

var (
	BatchSize     int
	BatchInterval time.Duration
)

func init() {
	var err error
	BatchSize, err = strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil || BatchSize == 0 {
		BatchSize = DEF_BATCH_SIZE
	}

	bi, err := strconv.Atoi(os.Getenv("BATCH_INTERVAL"))
	if err != nil || bi == 0 {
		BatchInterval = DEF_BATCH_INTERVAL
	} else {
		BatchInterval = time.Duration(bi) * time.Second
	}
}

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
	l := log.Logger{}

	lg, err := parseLog(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.Log(log.INFO, "cache write at %v \n", time.Now())
	bs, err := rh.cache.Write(ctx, lg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bs >= BatchSize {
		rh.BatchSizeChan <- struct{}{}
	}

	w.WriteHeader(http.StatusAccepted)
}

// HandleWebhookEvents handle the webhooks based on requirents
func HandleWebhookEvents(rh *receiverHandlers) {
	l := log.Logger{}

	for {
		t := time.NewTimer(BatchInterval)
		select {
		case <-rh.BatchSizeChan:
			l.Log(log.INFO, "cache batch size threshold crossed")
			logs, err := rh.cache.Evict()
			if err == nil {
				if err := client.RetryTimeout(logs, rh.client.HandleHook); err != nil {
					l.Log(log.INFO, "failed to release events to upstream server: %v", err)
				}
			}
		case <-t.C:
			l.Log(log.INFO, "cache batch interval crossed")
			logs, err := rh.cache.Evict()
			if err == nil {
				l.Log(log.INFO, "current batch size: %d", len(logs))
				// if empty cache do nothing
				if len(logs) <= 0 {
					continue
				}
				if err := client.RetryTimeout(logs, rh.client.HandleHook); err != nil {
					l.Log(log.ERROR, "failed to release events to upstream server: %v", err)
				}
			}
		}
	}
}

// -- for test ------
// HookHandler webhook handler function
func (rh *receiverHandlers) HooksHandler(w http.ResponseWriter, r *http.Request) {
	l := log.Logger{}

	logs := []model.Log{}

	//decode body
	err := json.NewDecoder(r.Body).Decode(&logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.Log(log.INFO, "hooks received: %v \n", logs)
	w.WriteHeader(http.StatusAccepted)
}
