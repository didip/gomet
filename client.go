package gomet

import (
	"bufio"
	"encoding/base64"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func NewClient(url string) (*Client, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		URL:          url,
		HTTPClient:   &http.Client{},
		HTTPRequest:  req,
		ResponseChan: make(chan []byte),
	}

	return c, nil
}

type Client struct {
	URL         string
	HTTPClient  *http.Client
	HTTPRequest *http.Request

	Retries         int
	MaxRetrySeconds int

	ResponseChan chan []byte

	OnConnectError      func(error)
	OnReadBytesError    func(error)
	OnBase64DecodeError func(error)
	OnPayloadReceived   func([]byte)

	mtx sync.RWMutex
}

func (c *Client) RetryDuration() time.Duration {
	if c.GetMaxRetrySeconds() <= 0 {
		c.SetMaxRetrySeconds(10) // default
	}
	return time.Duration(rand.Intn(c.GetMaxRetrySeconds())) * time.Second
}

func (c *Client) SetMaxRetrySeconds(seconds int) {
	c.mtx.Lock()
	c.MaxRetrySeconds = seconds
	c.mtx.Unlock()
}

func (c *Client) GetMaxRetrySeconds() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.MaxRetrySeconds
}

func (c *Client) SetRetries(retries int) {
	c.mtx.Lock()
	c.Retries = retries
	c.mtx.Unlock()
}

func (c *Client) GetRetries() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.Retries
}

func (c *Client) SetOnConnectError(f func(error)) {
	c.mtx.Lock()
	c.OnConnectError = f
	c.mtx.Unlock()
}

func (c *Client) SetOnReadBytesError(f func(error)) {
	c.mtx.Lock()
	c.OnReadBytesError = f
	c.mtx.Unlock()
}

func (c *Client) SetOnBase64DecodeError(f func(error)) {
	c.mtx.Lock()
	c.OnBase64DecodeError = f
	c.mtx.Unlock()
}

func (c *Client) SetOnPayloadReceived(f func([]byte)) {
	c.mtx.Lock()
	c.OnPayloadReceived = f
	c.mtx.Unlock()
}

func (c *Client) Connect() error {
	if c.URL == "" {
		return nil
	}

	var resp *http.Response
	var err error

	resp, err = c.HTTPClient.Do(c.HTTPRequest)
	if err != nil {
		if c.OnConnectError != nil {
			c.OnConnectError(err)
		}
		return err
	}

	if resp != nil && resp.Body != nil {
		reader := bufio.NewReader(resp.Body)
		for {
			payloadEncodedBytes, err := reader.ReadBytes('\n')
			if err != nil {
				if c.OnReadBytesError != nil {
					c.OnReadBytesError(err)
				}
				return err
			}

			payloadBytes, err := base64.StdEncoding.DecodeString(string(payloadEncodedBytes))
			if err != nil {
				if c.OnBase64DecodeError != nil {
					c.OnBase64DecodeError(err)
				}
				continue
			}

			if c.OnPayloadReceived != nil {
				c.OnPayloadReceived(payloadBytes)
			}
			c.ResponseChan <- payloadBytes
		}
	}

	return nil
}

func (c *Client) ConnectForever() {
	if c.URL == "" {
		return
	}

	if c.GetRetries() > 0 {
		for i := 0; i < c.GetRetries(); i++ {
			c.Connect()
			time.Sleep(c.RetryDuration())
		}

	} else {
		for {
			c.Connect()
			time.Sleep(c.RetryDuration())
		}
	}
}
