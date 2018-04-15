package gomet

import (
	"testing"
	"time"
)

func Test_RetryDuration(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	if c.RetryDuration() < time.Second {
		t.Errorf("Failed to generate a random duration that is > 1 second")
	}
}

func Test_SetGetRetries(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	c.SetRetries(3)

	if c.GetRetries() != 3 {
		t.Errorf("Failed to set/get Retries")
	}
}

func Test_SetGetMaxRetrySeconds(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	c.SetMaxRetrySeconds(20)

	if c.GetMaxRetrySeconds() != 20 {
		t.Errorf("Failed to set/get MaxRetrySeconds")
	}
}

func Test_SetOnConnectError(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	c.SetOnConnectError(func(error) { println("hello") })

	if c.OnConnectError == nil {
		t.Errorf("Failed to set OnConnectError hook")
	}
}

func Test_SetOnReadBytesError(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	c.SetOnReadBytesError(func(error) { println("hello") })

	if c.OnReadBytesError == nil {
		t.Errorf("Failed to set OnReadBytesError hook")
	}
}

func Test_SetOnBase64DecodeError(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	c.SetOnBase64DecodeError(func(error) { println("hello") })

	if c.OnBase64DecodeError == nil {
		t.Errorf("Failed to set OnBase64DecodeError hook")
	}
}

func Test_SetOnPayloadReceived(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	c.SetOnPayloadReceived(func([]byte) { println("hello") })

	if c.OnPayloadReceived == nil {
		t.Errorf("Failed to set OnPayloadReceived hook")
	}
}

func Test_BadConnect(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/stream")

	err := c.Connect()
	if err == nil {
		t.Errorf("error should be returned on bad connection")
	}
}
