package client_http

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Timeout int
}

type ClientRequest struct {
	Method     string
	Protocol   string
	Host       string
	Port       int
	Path       string
	Headers    http.Header
	FormValues []byte
	Payload    []byte
}

type ClientResponse struct {
	Headers    http.Header `json:"headers"`
	StatusCode int         `json:"status_code"`
}

func NewClient(timeout int) *Client {
	return &Client{
		Timeout: timeout,
	}
}

func (c *Client) Request(req *ClientRequest) (*ClientResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()

	request_method := req.Method
	request_endpoint := fmt.Sprintf(
		"%s://%s:%d/%s",
		req.Protocol,
		req.Host,
		req.Port,
		req.Path,
	)

	request_body := bytes.NewBuffer(req.Payload)
	if req.FormValues != nil {
		request_body = bytes.NewBuffer(req.FormValues)
	}

	r, err := http.NewRequestWithContext(
		ctx,
		request_method,
		request_endpoint,
		request_body,
	)
	if err != nil {
		log.Fatalf("Fail to create the request: %v", err)
		return nil, err
	}

	for key, values := range req.Headers {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Fatalf("Max timeout reached: %v", err)
			return nil, err
		}
	}

	response := ClientResponse{
		Headers:    res.Header,
		StatusCode: res.StatusCode,
	}

	return &response, nil
}
