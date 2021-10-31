package receiver

import (
	"context"
)

// WebhookReceiver interface to cover behaviours of a webhook receiver
type WebhookReceiver interface {
	HandleHook(ctx context.Context) error
}

type receiver struct {
}

// NewReceiver constructor for webhook recever type
func NewReceiver() WebhookReceiver {
	return &receiver{}
}

// HandleHooks the method to handle business logic of the hook
func (r *receiver) HandleHook(ctx context.Context) error {

	return nil
}
