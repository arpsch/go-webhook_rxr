package imemc

import (
	"context"
	"sync"

	"github.com/arpsch/go-webhook_rxr/model"
)

// Cache is the immemory cache holding the log
type Cache struct {
	mu   *sync.Mutex
	Logs []model.Log
}

// NewLogCache constructor for the cache
func NewLogCache() *Cache {
	lc := Cache{
		mu: new(sync.Mutex),
	}
	return &lc
}

// Write Write into the cache
func (c *Cache) Write(ctx context.Context, log *model.Log) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Logs = append(c.Logs, *log)

	return len(c.Logs), nil
}

// Evict evict data from the cache
func (c *Cache) Evict() ([]model.Log, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logs := c.Logs

	//clean the cache data
	c.Logs = nil

	return logs, nil
}
