package imemc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arpsch/go-webhook_rxr/model"
)

// TODO: Read from environment variable
const (
	BATCHSIZE     = 3 // number of items in cache
	BATCHINTERVAL = 30 * time.Second
)

// Cache is the immemory cache holding the log
type Cache struct {
	mu   *sync.Mutex
	Logs []model.Log

	// signal for go routine
	BatchSizeChan chan struct{}
}

// NewLogCache constructor for the cache
func NewLogCache() *Cache {
	lc := Cache{
		BatchSizeChan: make(chan struct{}),
		mu:            new(sync.Mutex),
	}
	return &lc
}

// Write Write into the cache
func (c *Cache) Write(ctx context.Context, log *model.Log) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Logs = append(c.Logs, *log)

	// if the cache has crossed the batch size then evict
	if len(c.Logs) >= BATCHSIZE {
		c.BatchSizeChan <- struct{}{}
	}

	return nil
}

// Evict evict data from the cache
func (c *Cache) Evict() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: call post end point
	fmt.Printf("%v\n", c.Logs)

	//clean the cache data
	c.Logs = nil

	return nil
}

// HandleWebhookEvents handle the webhooks based on requirents
func HandleWebhookEvents(cache *Cache) {
	for {
		t := time.NewTimer(BATCHINTERVAL)
		select {
		case <-cache.BatchSizeChan:
			cache.Evict()
		case <-t.C:
			if len(cache.Logs) > 0 {
				cache.Evict()
			}
		}
	}
}
