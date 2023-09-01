package connectors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/microlib/simple"
	"github.com/redis/go-redis/v9"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	Http        *http.Client
	RedisClient *redis.Client
	Logger      *simple.Logger
}

func NewClientConnections(logger *simple.Logger) Clients {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Connectors{Http: httpClient, Logger: logger, RedisClient: redis}
}

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Meta(info string) string {
	return info
}

func (c *Connectors) Do(req *http.Request) (*http.Response, error) {
	return c.Http.Do(req)
}

func (c *Connectors) Publish(ctx context.Context, topic string, payload interface{}) error {
	return c.RedisClient.Publish(ctx, topic, payload).Err()
}
