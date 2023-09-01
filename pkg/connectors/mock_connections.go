package connectors

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/microlib/simple"
)

// Mock all connections
type MockConnectors struct {
	Http   *http.Client
	Logger *simple.Logger
	Flag   string
}

func (c *MockConnectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Meta(flag string) string {
	c.Flag = flag
	return flag
}

func (c *MockConnectors) Do(req *http.Request) (*http.Response, error) {
	if c.Flag == "true" {
		return nil, errors.New("forced http error")
	}
	return c.Http.Do(req)
}

func (c *MockConnectors) Publish(ctx context.Context, topic string, payload interface{}) error {
	return nil
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestConnectors(file string, code int, logger *simple.Logger) Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	var data []byte
	var err error
	if len(file) > 0 {
		data, err = os.ReadFile(file)
		if err != nil {
			logger.Error(fmt.Sprintf("file data %v\n", err))
			panic(err)
		}
	} else {
		data = []byte("")
	}
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: io.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &MockConnectors{Http: httpclient, Logger: logger, Flag: "false"}
	return conns
}
