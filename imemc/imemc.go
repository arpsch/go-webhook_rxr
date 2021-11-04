package imemc

import (
	"context"
	"os"
	"strconv"

	"github.com/arpsch/go-webhook_rxr/client"
	"github.com/arpsch/go-webhook_rxr/model"
)

const (
	DEF_BATCH_SIZE = 3
)

var (
	BatchSize int
)

func init() {
	var err error
	BatchSize, err = strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil || BatchSize == 0 {
		BatchSize = DEF_BATCH_SIZE
	}
}

// Cache is the immemory cache holding the log
type Cache struct {
	// cache log logs
	Logs []model.Log

	// channel to manage insertions and eviction
	ch chan model.Log
}

// NewLogCache constructor for the cache
func NewLogCache() *Cache {
	lc := Cache{
		ch: make(chan model.Log),
	}
	return &lc
}

// Write Write into the cache
func (c *Cache) Write(ctx context.Context, log *model.Log) error {

	c.ch <- *log

	return nil
}

// Evict evict data from the cache
func (c *Cache) Evict() {
	// singal the eviction
	// TODO: for some reason, emtpy struct comparison is not working, so using Title field to evict
	// c.ch  <- model.Log{} // send empty Log struct to signal eviction
	c.ch <- model.Log{
		Title: "evict",
	}
}

// CacheController control the cache population/eviction based on the rule
// It will call the handler supplied when eviction is signalled
func CacheController(c *Cache, cl *client.Client, fn func([]model.Log, *client.Client) error) {

	for {

		l := <-c.ch

		// dont cache if this Log is signaling to evict
		if l.Title != "evict" {
			c.Logs = append(c.Logs, l)
		}

		// TODO: for some reason, emtpy struct comparison is not wotking, so using Title field to evict
		//if len(c.Logs) > 4 || l == (model.Log{}) {
		if len(c.Logs) >= BatchSize || l.Title == "evict" {
			// call the eviction handler
			go fn(c.Logs, cl)
			//reset the logs
			c.Logs = nil
		}
	}
}
