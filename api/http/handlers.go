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
	DEF_BATCH_INTERVAL = 30 * time.Second
)

var (
	BatchInterval time.Duration
)

func init() {
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
}

// NewReceiverHandlers constructor for Receiver
func NewReceiverHandlers(cache *imemc.Cache, client *client.Client) *receiverHandlers {
	return &receiverHandlers{
		cache:  cache,
		client: client,
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
	err = rh.cache.Write(ctx, lg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// HandleWebhookEvents handle the webhooks based on requirents
func HandleBatchIntervalTimeout(rh *receiverHandlers) {
	l := log.Logger{}

	for {
		t := time.NewTimer(BatchInterval)
		select {
		case <-t.C:
			l.Log(log.INFO, "cache batch interval crossed")
			rh.cache.Evict()
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
