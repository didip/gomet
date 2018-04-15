package gomet

import (
	"testing"
)

func Test_SetHTTPResponseHeaders(t *testing.T) {
	b := NewBroadcaster()
	b.SetHTTPResponseHeaders(map[string]string{
		"foo": "bar",
	})

	if b.HTTPResponseHeaders["foo"] != "bar" {
		t.Errorf("Failed to set custom HTTP header: %v", b.HTTPResponseHeaders["foo"])
	}
}
