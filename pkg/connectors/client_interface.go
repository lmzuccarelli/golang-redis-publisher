package connectors

import (
	"net/http"

	"context"
)

// Client Interface - used as a receiver and can be overridden for testing
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Publish(ctx context.Context, topic string, payload interface{}) error
	Do(req *http.Request) (*http.Response, error)
}
