package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestClientDataAttachTimeoutClosesWith1013(t *testing.T) {
	reg := newRegistry(64)
	handler := &relayHandler{
		reg:               reg,
		dataAttachTimeout: 50 * time.Millisecond,
	}
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?serverId=test-server&role=client&v=2"
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial client: %v", err)
	}
	defer c.Close()

	c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err = c.ReadMessage()
	if err == nil {
		t.Fatal("expected connection to be closed due to attach timeout")
	}
	closeErr, ok := err.(*websocket.CloseError)
	if !ok {
		t.Fatalf("expected CloseError, got: %v", err)
	}
	if closeErr.Code != 1013 {
		t.Fatalf("expected close code 1013, got %d", closeErr.Code)
	}
	if !strings.Contains(closeErr.Text, "Data route unavailable") {
		t.Fatalf("expected reason 'Data route unavailable', got %q", closeErr.Text)
	}
}
