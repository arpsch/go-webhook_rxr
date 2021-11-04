package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	log "github.com/arpsch/go-webhook_rxr/logger"
	"github.com/arpsch/go-webhook_rxr/model"
)

const (
	RetryInterval = 2 * time.Second
	RetryCount    = 3
)

var (
	endpoint string
)

type Client struct {
	Endpoint string
}

func init() {
	endpoint = os.Getenv("ENDPOINT")
	// set default value is not set through ENV
	if endpoint == "" {
		endpoint = "http://127.0.0.1:9999/logs"
	}
}

// NewReceiver constructor for webhook recever type
func NewClient() *Client {
	return &Client{
		// TODO: to be read from environment
		Endpoint: endpoint,
	}
}

// HandleHooks the method to handle business logic of the hook
func (c *Client) PostHookEndpoint(ctx context.Context, logs []model.Log) (int, time.Duration, error) {
	json_data, err := json.Marshal(logs)
	if err != nil {
		return 0, 0, err
	}

	s := time.Now()
	resp, err := http.Post(c.Endpoint, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return 0, 0, err
	}
	e := time.Now()

	if resp.StatusCode != http.StatusAccepted {
		return 0, 0, errors.New("request failed")
	}

	return resp.StatusCode, e.Sub(s), nil
}

// RetryTimeout calls the client endpoint with retry constraint
// Attempt 3 times at an interval of 2 seconds each
func InvokePostEnpoint(logs []model.Log, c *Client) error {
	l := log.Logger{}

	// if logs are empty, do nothing
	if len(logs) == 0 {
		l.Log(log.WARN, "empty cache, nothing to do")
		return errors.New("empty cache, nothing to do")
	}

	//set up context with time out for upstream request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < RetryCount; i++ {
		if sc, d, err := c.PostHookEndpoint(ctx, logs); err == nil {
			l.Log(log.INFO, "Posted batch of size: %d in %v seconds with status: %d \n", len(logs), d, sc)
			return nil
		}

		l.Log(log.WARN, "wait %s before trying again\n", RetryInterval)
		t := time.NewTimer(RetryInterval)

		select {
		case <-ctx.Done():
			l.Log(log.WARN, "time expired 2 : %v\n", ctx.Err())
			t.Stop()
			return errors.New("time expired")
		case <-t.C:
			l.Log(log.WARN, "retry again -  count %d\n", i)
		}
	}
	return errors.New("failed to update hooks")
}
