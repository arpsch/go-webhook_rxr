package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arpsch/go-webhook_rxr/model"
)

const (
	RetryInterval = 2 * time.Second
	RetryCount    = 3
)

type Client struct {
	Endpoint string
}

// NewReceiver constructor for webhook recever type
func NewClient() *Client {
	return &Client{
		// TODO: to be read from environment
		Endpoint: "http://requestbin.net",
	}
}

// HandleHooks the method to handle business logic of the hook
func (c *Client) HandleHook(ctx context.Context, logs []model.Log) error {
	json_data, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	_, err = http.Post(c.Endpoint, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	return nil
}

// RetryTimeout calls the client endpoint with retry constraint
// Attempt 3 times at an interval of 2 seconds each
func RetryTimeout(logs []model.Log,
	check func(context.Context, []model.Log) error) error {
	//set up context for upstream request
	ctx := context.Background()

	for i := 0; i < RetryCount; i++ {
		if err := check(ctx, logs); err == nil {
			log.Printf("finished successfully in attempt: %d\n", i)
			return err
		}
		if ctx.Err() != nil {
			log.Printf("time expired 1 : %v\n", ctx.Err())
			return errors.New(ctx.Err().Error())
		}
		log.Printf("wait %s before trying again\n", RetryInterval)
		t := time.NewTimer(RetryInterval)
		select {
		case <-ctx.Done():
			log.Printf("time expired 2 : %v\n", ctx.Err())
			t.Stop()
			return errors.New("time expired")
		case <-t.C:
			log.Printf("retry again -  count %d\n", i)
		}
	}
	return nil
}
